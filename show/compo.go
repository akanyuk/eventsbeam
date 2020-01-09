package show

import (
	"TEST-LOCAL/eventsbeam/show/config"
	"TEST-LOCAL/eventsbeam/show/storage"
	"TEST-LOCAL/eventsbeam/web/response"

	"errors"
	"path/filepath"
	"sync"
)

const compoConfigFileName = "compo.yaml"

type compo struct {
	sync.RWMutex
	compos     []config.Compo
	configPath string
}

type Comper interface {
	Init() error
	Validate(config.Compo, config.Compo) []response.ErrorItem
	Compos() []config.Compo
	Create(config.Compo) error
	Read(string) (config.Compo, error)
	Update(alias string, compo config.Compo) error
	Delete(alias string) error
}

func NewComper(basePath string) Comper {
	return &compo{
		configPath: filepath.Join(basePath, compoConfigFileName),
	}
}

func (c *compo) Init() error {
	compos, err := storage.LoadCompos(c.configPath)
	if err != nil {
		compos = []config.Compo{}
	}

	c.Lock()
	defer c.Unlock()

	c.compos = compos

	return nil
}

func (c *compo) Validate(compo config.Compo, oldCompo config.Compo) []response.ErrorItem {
	errorItems := make([]response.ErrorItem, 0)

	if compo.Alias == "" {
		errorItems = append(errorItems, response.ErrorItem{Code: "alias", Message: "alias can not be empty"})
	} else if compo.Alias != oldCompo.Alias {
		_, exist := c.getCompo(compo.Alias)
		if exist {
			errorItems = append(errorItems, response.ErrorItem{Code: "alias", Message: "alias already exists"})
		}
	}

	if compo.Title == "" {
		errorItems = append(errorItems, response.ErrorItem{Code: "title", Message: "title can not be empty"})
	}

	return errorItems
}

func (c *compo) Compos() []config.Compo {
	return c.compos
}

func (c *compo) Read(alias string) (config.Compo, error) {
	compo, exist := c.getCompo(alias)
	if !exist {
		return config.Compo{}, errors.New("compo not found")
	}

	return compo, nil
}

func (c *compo) Create(compo config.Compo) error {
	_, exist := c.getCompo(compo.Alias)
	if exist {
		return errors.New("alias already exist")
	}

	c.Lock()
	defer c.Unlock()

	c.compos = append(c.compos, compo)
	if err := storage.SaveCompos(c.compos, c.configPath); err != nil {
		return err
	}

	return nil
}

func (c *compo) Update(alias string, updatedCompo config.Compo) error {
	_, exist := c.getCompo(alias)
	if !exist {
		return errors.New("not found")
	}

	c.Lock()
	defer c.Unlock()

	for key, item := range c.compos {
		if item.Alias == alias {
			c.compos[key] = updatedCompo
		}
	}

	if err := storage.SaveCompos(c.compos, c.configPath); err != nil {
		return err
	}

	return nil
}

func (c *compo) Delete(alias string) error {
	_, exist := c.getCompo(alias)
	if !exist {
		return errors.New("not found")
	}

	c.Lock()
	defer c.Unlock()

	for key, item := range c.compos {
		if item.Alias == alias {
			c.compos = append(c.compos[:key], c.compos[key+1:]...)
		}
	}

	if err := storage.SaveCompos(c.compos, c.configPath); err != nil {
		return err
	}

	return nil
}

func (c *compo) getCompo(alias string) (config.Compo, bool) {
	c.Lock()
	defer c.Unlock()

	for _, compo := range c.compos {
		if compo.Alias == alias {
			return compo, true
		}
	}

	return config.Compo{}, false
}
