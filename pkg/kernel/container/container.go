package container

import (
    "sync"
    "reflect"
    "fmt"
    "strings"
    "errors"
)

var (
    ErrFactoryNotFound = errors.New("factory not found")
)

type factory = func() (interface{}, error)

// 容器
type Container struct {
    sync.Mutex
    singletons map[string]interface{}
    factories  map[string]factory
}

// 容器实例化
func NewContainer() *Container {
    return &Container{
        singletons: make(map[string]interface{}),
        factories:  make(map[string]factory),
    }
}

// 注册单例对象
func (p *Container) SetSingleton(name string, singleton interface{}) {
    p.Lock()
    p.singletons[name] = singleton
    p.Unlock()
}

// 获取单例对象
func (p *Container) GetSingleton(name string) interface{} {
    return p.singletons[name]
}

// 获取实例对象
func (p *Container) GetPrototype(name string) (interface{}, error) {
    factory, ok := p.factories[name]
    if !ok {
        return nil, ErrFactoryNotFound
    }
    return factory()
}

// 设置实例对象工厂
func (p *Container) SetPrototype(name string, factory factory) {
    p.Lock()
    p.factories[name] = factory
    p.Unlock()
}

// 注入依赖
func (p *Container) Ensure(instance interface{}) error {
    elemType := reflect.TypeOf(instance).Elem()
    ele := reflect.ValueOf(instance).Elem()
    for i := 0; i < elemType.NumField(); i++ { // 遍历字段
        fieldType := elemType.Field(i)
        tag := fieldType.Tag.Get("di") // 获取tag
        diName := p.injectName(tag)
        if diName == "" {
            continue
        }
        var (
            diInstance interface{}
            err        error
        )
        if p.isSingleton(tag) {
            diInstance = p.GetSingleton(diName)
        }
        if p.isPrototype(tag) {
            diInstance, err = p.GetPrototype(diName)
        }
        if err != nil {
            return err
        }
        if diInstance == nil {
            return errors.New(diName + " dependency not found")
        }
        ele.Field(i).Set(reflect.ValueOf(diInstance))
    }
    return nil
}

// 获取需要注入的依赖名称
func (p *Container) injectName(tag string) string {
    tags := strings.Split(tag, ",")
    if len(tags) == 0 {
        return ""
    }
    return tags[0]
}

// 检测是否单例依赖
func (p *Container) isSingleton(tag string) bool {
    tags := strings.Split(tag, ",")
    for _, name := range tags {
        if name == "prototype" {
            return false
        }
    }
    return true
}

// 检测是否实例依赖
func (p *Container) isPrototype(tag string) bool {
    tags := strings.Split(tag, ",")
    for _, name := range tags {
        if name == "prototype" {
            return true
        }
    }
    return false
}

// 打印容器内部实例
func (p *Container) String() string {
	lines := make([]string, 0, len(p.singletons)+len(p.factories)+2)
	lines = append(lines, "singletons:")
	for name, item := range p.singletons {
		if item == nil {
			line := fmt.Sprintf("  %s: %s %s", name, "<nil>", "<nil>")
			lines = append(lines, line)
			continue
		}

		line := fmt.Sprintf("  %s: %p %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	lines = append(lines, "factories:")
	for name, item := range p.factories {
		if item == nil {
			line := fmt.Sprintf("  %s: %s %s", name, "<nil>", "<nil>")
			lines = append(lines, line)
			continue
		}

		line := fmt.Sprintf("  %s: %p %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
