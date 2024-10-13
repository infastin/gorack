package openapi

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"gopkg.in/yaml.v3"
)

type item struct {
	key  string
	val  *yaml.Node
	next *item
}

func newItem(key string, value *yaml.Node) *item {
	return &item{
		key:  key,
		val:  value,
		next: nil,
	}
}

type openAPICommon struct {
	version  *item
	info     *item
	security *item
	first    *item
	last     *item
}

type openAPI interface {
	yaml.Marshaler
	setBasePath(basePath string) error
	setBearerAuth(bearerFormat string) error
}

func parseOpenAPI(node *yaml.Node) (openAPI, error) {
	if node.Kind != yaml.DocumentNode {
		return nil, errors.New(`invalid OpenAPI spec: expected object`)
	}
	node = node.Content[0]

	var common openAPICommon
	var otherFields []*item
	var prev *item

	for i := 0; i < len(node.Content); i += 2 {
		cur := newItem(node.Content[i].Value, node.Content[i+1])

		if prev != nil {
			prev.next = cur
		}
		prev = cur

		if common.first == nil {
			common.first = cur
		}
		common.last = cur

		switch cur.key {
		case "openapi", "swagger":
			if cur.val.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf(`invalid OpenAPI spec: %q is not a string`, cur.key)
			}
			common.version = cur
		case "info":
			if cur.val.Kind != yaml.MappingNode {
				return nil, errors.New(`invalid OpenAPI spec: "info" is not an object`)
			}
			common.info = cur
		case "security":
			if cur.val.Kind != yaml.SequenceNode {
				return nil, errors.New(`invalid OpenAPI spec: "security" is not an array`)
			}
			common.security = cur
		default:
			otherFields = append(otherFields, cur)
		}
	}

	if common.version == nil {
		return nil, errors.New(`invalid OpenAPI spec: "openapi" or "swagger" not specified`)
	}

	if common.info == nil {
		return nil, errors.New(`invalid OpenAPI spec: "info" not specified`)
	}

	switch common.version.val.Value {
	case "2.0":
		v2 := &openAPIv2{openAPICommon: common}
		for _, kv := range otherFields {
			switch kv.key {
			case "basePath":
				if kv.val.Kind != yaml.ScalarNode {
					return nil, errors.New(`invalid OpenAPI spec: "basePath" is not a string`)
				}
				v2.basePath = kv
			case "securityDefinitions":
				if kv.val.Kind != yaml.MappingNode {
					return nil, errors.New(`invalid OpenAPI spec: "securityDefinitions" is not an object`)
				}
				v2.securityDefinitions = kv
			default:
				v2.otherFields = append(v2.otherFields, kv)
			}
		}
		return v2, nil
	case "3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.1.0":
		v3 := &openAPIv3{openAPICommon: common}
		for _, kv := range otherFields {
			switch kv.key {
			case "servers":
				if kv.val.Kind != yaml.SequenceNode {
					return nil, errors.New(`invalid OpenAPI spec: "servers" is not an array`)
				}
				v3.servers = kv
			case "components":
				if kv.val.Kind != yaml.MappingNode {
					return nil, errors.New(`invalid OpenAPI spec: "components" is not an object`)
				}
				v3.components = kv
			default:
				v3.otherFields = append(v3.otherFields, kv)
			}
		}
		return v3, nil
	default:
		return nil, fmt.Errorf("invalid OpenAPI spec: invalid version %q", common.version.val.Value)
	}
}

type openAPIv2 struct {
	openAPICommon
	basePath            *item
	securityDefinitions *item
	otherFields         []*item
}

func (o *openAPIv2) setBasePath(basePath string) error {
	if o.basePath == nil {
		o.basePath = newItem("basePath", nil)
		o.basePath.next = o.info.next
		o.info.next = o.basePath
	}
	o.basePath.val = &yaml.Node{Kind: yaml.ScalarNode, Value: path.Clean(basePath)}
	return nil
}

func (o *openAPIv2) setBearerAuth(format string) error {
	if o.securityDefinitions == nil {
		o.securityDefinitions = newItem("securityDefinitions", &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: nil,
		})
		o.last.next = o.securityDefinitions
		o.last = o.securityDefinitions
	}

	o.securityDefinitions.val.Content = append(o.securityDefinitions.val.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "BearerAuth"},
		&yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "type"}, {Kind: yaml.ScalarNode, Value: "apiKey"},
			{Kind: yaml.ScalarNode, Value: "in"}, {Kind: yaml.ScalarNode, Value: "header"},
			{Kind: yaml.ScalarNode, Value: "name"}, {Kind: yaml.ScalarNode, Value: "Authorization"},
		}},
	)

	if o.security == nil {
		o.security = newItem("security", &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: nil,
		})
		o.security.next = o.info.next
		o.info.next = o.security
	}

	o.security.val.Content = append(o.security.val.Content, &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "BearerAuth"},
			{Kind: yaml.SequenceNode, Content: nil},
		},
	})

	return nil
}

func (o *openAPIv2) MarshalYAML() (any, error) {
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, 0),
	}

	for cur := o.first; cur != nil; cur = cur.next {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: cur.key},
			cur.val,
		)
	}

	return node, nil
}

type openAPIv3 struct {
	openAPICommon
	servers     *item
	components  *item
	otherFields []*item
}

func (o *openAPIv3) setBasePath(basePath string) error {
	if o.servers == nil {
		o.servers = newItem("servers", &yaml.Node{
			Kind: yaml.SequenceNode, Content: nil,
		})
		o.servers.next = o.info.next
		o.info.next = o.servers
	}

	if len(o.servers.val.Content) == 0 {
		o.servers.val.Content = append(o.servers.val.Content, &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "url"},
				{Kind: yaml.ScalarNode, Value: basePath},
			},
		})
	} else {
		for _, server := range o.servers.val.Content {
			if server.Kind != yaml.MappingNode {
				return errors.New(`invalid OpenAPI spec: "servers" is not an array of objects`)
			}
			for i := 0; i < len(server.Content); i += 2 {
				k, v := server.Content[i], server.Content[i+1]
				if k.Value == "url" {
					if v.Kind != yaml.ScalarNode {
						return fmt.Errorf(`invalid OpenAPI spec: "servers.[%d].url" is not a string`, i/2)
					}
					v.Value = basePath
				}
			}
		}
	}

	return nil
}

func (o *openAPIv3) setBearerAuth(format string) error {
	if o.components == nil {
		o.components = newItem("components", &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: nil,
		})
		o.last.next = o.components
		o.last = o.components
	}

	var securitySchemes *yaml.Node
	for i := 0; i < len(o.components.val.Content); i += 2 {
		k, v := o.components.val.Content[i], o.components.val.Content[i+1]
		if k.Value == "securitySchemes" {
			if v.Kind != yaml.MappingNode {
				return errors.New(`invalid OpenAPI spec: "components/securitySchemes" is not an object`)
			}
			securitySchemes = v
			break
		}
	}

	if securitySchemes == nil {
		securitySchemes = &yaml.Node{
			Kind:    yaml.MappingNode,
			Content: make([]*yaml.Node, 0, 1),
		}
		o.components.val.Content = append(o.components.val.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "securitySchemes"},
			securitySchemes,
		)
	}

	securitySchemes.Content = append(securitySchemes.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "BearerAuth"},
		&yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "type"}, {Kind: yaml.ScalarNode, Value: "http"},
				{Kind: yaml.ScalarNode, Value: "scheme"}, {Kind: yaml.ScalarNode, Value: "bearer"},
				{Kind: yaml.ScalarNode, Value: "bearerFormat"}, {Kind: yaml.ScalarNode, Value: format},
			},
		},
	)

	if o.security == nil {
		o.security = newItem("security", &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: nil,
		})
		o.security.next = o.info.next
		o.info.next = o.security
	}

	o.security.val.Content = append(o.security.val.Content, &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "BearerAuth"},
			{Kind: yaml.SequenceNode, Content: nil},
		},
	})

	return nil
}

func (o *openAPIv3) MarshalYAML() (any, error) {
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, 0),
	}

	for cur := o.first; cur != nil; cur = cur.next {
		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: cur.key},
			cur.val,
		)
	}

	return node, nil
}

type AppendOptions struct {
	BasePath   string
	BearerAuth string
}

func Append(spec []byte, opts *AppendOptions) ([]byte, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(spec, &document); err != nil {
		return nil, err
	}

	openapi, err := parseOpenAPI(&document)
	if err != nil {
		return nil, err
	}

	if opts.BasePath != "" {
		if err := openapi.setBasePath(opts.BasePath); err != nil {
			return nil, err
		}
	}

	if opts.BearerAuth != "" {
		if err := openapi.setBearerAuth(opts.BearerAuth); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(openapi); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
