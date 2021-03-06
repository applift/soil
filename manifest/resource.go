package manifest

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/copystructure"
)

const (
	openResourcePrefix    = "resource"
	resourceRequestPrefix = "resource.request"
)

// Resources are referenced by ${resource.<kind>.<pod>.name}
type Resource struct {

	// Resource name unique within pod
	Name string `hcl:"-"`

	// Resource type
	Kind string `hcl:"-"`

	// Add "resource.<Type>.<PodName>.<Name>" = "true" to pod allocation constraint
	Required bool

	// Request config
	Config map[string]interface{} `hcl:"-"`
}

func defaultResource() (r Resource) {
	r = Resource{
		Required: true,
	}
	return
}

func (r Resource) Clone() (res Resource) {
	res1, _ := copystructure.Copy(r)
	res = res1.(Resource)
	return
}

// GetID resource ID
func (r *Resource) GetID(podName string) (res string) {
	res = fmt.Sprintf("%s.%s", podName, r.Name)
	return
}

// Returns "resource.request.<kind>.allow": "true"
func (r *Resource) GetRequestConstraint() (res Constraint) {
	res = Constraint{
		fmt.Sprintf("${%s.%s.allow}", resourceRequestPrefix, r.Kind): "true",
	}
	return
}

// Returns required constraint for provision with allocated resource
func (r *Resource) GetAllocationConstraint(podName string) (res Constraint) {
	res = Constraint{}
	if r.Required {
		res[fmt.Sprintf("${%s.%s.%s.allocated}", openResourcePrefix, r.Kind, r.GetID(podName))] = "true"
	}
	return
}

// Returns `resource.<kind>.<pod>.<name>.__values_json`
func (r *Resource) GetValuesKey(podName string) (res string) {
	res = fmt.Sprintf("%s.%s.%s.__values", openResourcePrefix, r.Kind, r.GetID(podName))
	return
}

func (r *Resource) parseAst(raw *ast.ObjectItem) (err error) {
	if len(raw.Keys) != 2 {
		err = fmt.Errorf(`resource should be "type" "name"`)
		return
	}
	r.Kind = raw.Keys[0].Token.Value().(string)
	r.Name = raw.Keys[1].Token.Value().(string)
	if err = hcl.DecodeObject(r, raw); err != nil {
		return
	}
	if err = hcl.DecodeObject(&r.Config, raw.Val); err != nil {
		return
	}
	delete(r.Config, "required")
	delete(r.Config, "name")
	delete(r.Config, "kind")
	return
}
