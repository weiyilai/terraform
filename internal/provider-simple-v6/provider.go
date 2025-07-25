// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// simple provider a minimal provider implementation for testing
package simple

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/hashicorp/terraform/internal/configs/configschema"
	"github.com/hashicorp/terraform/internal/providers"
)

type simple struct {
	schema providers.GetProviderSchemaResponse
}

func Provider() providers.Interface {
	simpleResource := providers.Schema{
		Body: &configschema.Block{
			Attributes: map[string]*configschema.Attribute{
				"id": {
					Computed: true,
					Type:     cty.String,
				},
				"value": {
					Optional: true,
					Type:     cty.String,
				},
			},
		},
		Identity: &configschema.Object{
			Attributes: map[string]*configschema.Attribute{
				"id": {
					Type:     cty.String,
					Required: true,
				},
			},
			Nesting: configschema.NestingSingle,
		},
	}

	return simple{
		schema: providers.GetProviderSchemaResponse{
			Provider: providers.Schema{
				Body: nil,
			},
			ResourceTypes: map[string]providers.Schema{
				"simple_resource": simpleResource,
			},
			DataSources: map[string]providers.Schema{
				"simple_resource": simpleResource,
			},
			EphemeralResourceTypes: map[string]providers.Schema{
				"simple_resource": simpleResource,
			},
			ListResourceTypes: map[string]providers.Schema{
				"simple_resource": {
					Body: &configschema.Block{
						Attributes: map[string]*configschema.Attribute{
							"value": {
								Optional: true,
								Type:     cty.String,
							},
						},
					},
				},
			},
			Actions: map[string]providers.ActionSchema{},
			ServerCapabilities: providers.ServerCapabilities{
				PlanDestroy:               true,
				GetProviderSchemaOptional: true,
			},
			Functions: map[string]providers.FunctionDecl{
				"noop": {
					Parameters: []providers.FunctionParam{
						{
							Name:               "noop",
							Type:               cty.DynamicPseudoType,
							AllowNullValue:     true,
							AllowUnknownValues: true,
							Description:        "any value",
							DescriptionKind:    configschema.StringPlain,
						},
					},
					ReturnType:      cty.DynamicPseudoType,
					Description:     "noop takes any single argument and returns the same value",
					DescriptionKind: configschema.StringPlain,
				},
			},
		},
	}
}

func (s simple) GetProviderSchema() providers.GetProviderSchemaResponse {
	return s.schema
}

func (s simple) GetResourceIdentitySchemas() providers.GetResourceIdentitySchemasResponse {
	return providers.GetResourceIdentitySchemasResponse{
		IdentityTypes: map[string]providers.IdentitySchema{
			"simple_resource": {
				Version: 0,
				Body: &configschema.Object{
					Attributes: map[string]*configschema.Attribute{
						"id": {
							Type:     cty.String,
							Required: true,
						},
					},
					Nesting: configschema.NestingSingle,
				},
			},
		},
	}
}

func (s simple) ValidateProviderConfig(req providers.ValidateProviderConfigRequest) (resp providers.ValidateProviderConfigResponse) {
	return resp
}

func (s simple) ValidateResourceConfig(req providers.ValidateResourceConfigRequest) (resp providers.ValidateResourceConfigResponse) {
	return resp
}

func (s simple) ValidateDataResourceConfig(req providers.ValidateDataResourceConfigRequest) (resp providers.ValidateDataResourceConfigResponse) {
	return resp
}

func (s simple) ValidateListResourceConfig(req providers.ValidateListResourceConfigRequest) (resp providers.ValidateListResourceConfigResponse) {
	return resp
}

func (p simple) UpgradeResourceState(req providers.UpgradeResourceStateRequest) (resp providers.UpgradeResourceStateResponse) {
	ty := p.schema.ResourceTypes[req.TypeName].Body.ImpliedType()
	val, err := ctyjson.Unmarshal(req.RawStateJSON, ty)
	resp.Diagnostics = resp.Diagnostics.Append(err)
	resp.UpgradedState = val
	return resp
}

func (p simple) UpgradeResourceIdentity(req providers.UpgradeResourceIdentityRequest) (resp providers.UpgradeResourceIdentityResponse) {
	schema := p.GetResourceIdentitySchemas().IdentityTypes[req.TypeName].Body
	ty := schema.ImpliedType()
	val, err := ctyjson.Unmarshal(req.RawIdentityJSON, ty)
	resp.Diagnostics = resp.Diagnostics.Append(err)
	resp.UpgradedIdentity = val
	return resp
}

func (s simple) ConfigureProvider(providers.ConfigureProviderRequest) (resp providers.ConfigureProviderResponse) {
	return resp
}

func (s simple) Stop() error {
	return nil
}

func (s simple) ReadResource(req providers.ReadResourceRequest) (resp providers.ReadResourceResponse) {
	// just return the same state we received
	resp.NewState = req.PriorState
	resp.Identity = req.CurrentIdentity
	return resp
}

func (s simple) PlanResourceChange(req providers.PlanResourceChangeRequest) (resp providers.PlanResourceChangeResponse) {
	if req.ProposedNewState.IsNull() {
		// destroy op
		resp.PlannedState = req.ProposedNewState

		// signal that this resource was properly planned for destruction,
		// verifying that the schema capabilities with PlanDestroy took effect.
		resp.PlannedPrivate = []byte("destroy planned")
		return resp
	}

	m := req.ProposedNewState.AsValueMap()
	_, ok := m["id"]
	if !ok {
		m["id"] = cty.UnknownVal(cty.String)
	}

	resp.PlannedState = cty.ObjectVal(m)
	return resp
}

func (s simple) ApplyResourceChange(req providers.ApplyResourceChangeRequest) (resp providers.ApplyResourceChangeResponse) {
	if req.PlannedState.IsNull() {
		// make sure this was transferred from the plan action
		if string(req.PlannedPrivate) != "destroy planned" {
			resp.Diagnostics = resp.Diagnostics.Append(fmt.Errorf("resource not planned for destroy, private data %q", req.PlannedPrivate))
		}

		resp.NewState = req.PlannedState
		resp.NewIdentity = req.PlannedIdentity
		return resp
	}

	m := req.PlannedState.AsValueMap()
	_, ok := m["id"]
	if !ok {
		m["id"] = cty.StringVal(time.Now().String())
	}
	resp.NewState = cty.ObjectVal(m)

	return resp
}

func (s simple) ImportResourceState(providers.ImportResourceStateRequest) (resp providers.ImportResourceStateResponse) {
	resp.Diagnostics = resp.Diagnostics.Append(errors.New("unsupported"))
	return resp
}

func (s simple) MoveResourceState(providers.MoveResourceStateRequest) (resp providers.MoveResourceStateResponse) {
	// We don't expose the move_resource_state capability, so this should never
	// be called.
	resp.Diagnostics = resp.Diagnostics.Append(errors.New("unsupported"))
	return resp
}

func (s simple) ReadDataSource(req providers.ReadDataSourceRequest) (resp providers.ReadDataSourceResponse) {
	m := req.Config.AsValueMap()
	m["id"] = cty.StringVal("static_id")
	resp.State = cty.ObjectVal(m)
	return resp
}

func (p simple) ValidateEphemeralResourceConfig(req providers.ValidateEphemeralResourceConfigRequest) (resp providers.ValidateEphemeralResourceConfigResponse) {
	return resp
}

func (s simple) OpenEphemeralResource(req providers.OpenEphemeralResourceRequest) (resp providers.OpenEphemeralResourceResponse) {
	// we only have one type, so no need to check
	m := req.Config.AsValueMap()
	m["id"] = cty.StringVal("ephemeral secret")
	resp.Result = cty.ObjectVal(m)
	resp.Private = []byte("private data")
	resp.RenewAt = time.Now().Add(time.Second)
	return resp
}

func (s simple) RenewEphemeralResource(req providers.RenewEphemeralResourceRequest) (resp providers.RenewEphemeralResourceResponse) {
	log.Printf("[DEBUG] renewing ephemeral resource")
	if string(req.Private) != "private data" {
		resp.Diagnostics = resp.Diagnostics.Append(fmt.Errorf("invalid private data %q, cannot renew ephemeral resource", req.Private))
	}
	resp.Private = req.Private
	resp.RenewAt = time.Now().Add(time.Second)
	return resp
}

func (s simple) CloseEphemeralResource(req providers.CloseEphemeralResourceRequest) (resp providers.CloseEphemeralResourceResponse) {
	log.Printf("[DEBUG] closing ephemeral resource")
	if string(req.Private) != "private data" {
		resp.Diagnostics = resp.Diagnostics.Append(fmt.Errorf("invalid private data %q, cannot close ephemeral resource", req.Private))
	}
	return resp
}

func (s simple) CallFunction(req providers.CallFunctionRequest) (resp providers.CallFunctionResponse) {
	if req.FunctionName != "noop" {
		resp.Err = fmt.Errorf("CallFunction for undefined function %q", req.FunctionName)
		return resp
	}

	resp.Result = req.Arguments[0]
	return resp
}

func (s simple) ListResource(req providers.ListResourceRequest) (resp providers.ListResourceResponse) {
	vals := make([]cty.Value, 0)

	staticVal := cty.StringVal("static_value")
	m := req.Config.AsValueMap()
	if val, ok := m["value"]; ok && val != cty.NilVal {
		staticVal = val
	}

	obj := map[string]cty.Value{
		"display_name": cty.StringVal("static_display_name"),
		"identity": cty.ObjectVal(map[string]cty.Value{
			"id": cty.StringVal("static_id"),
		}),
	}
	if req.IncludeResourceObject {
		obj["state"] = cty.ObjectVal(map[string]cty.Value{
			"id":    cty.StringVal("static_id"),
			"value": staticVal,
		})
	}
	vals = append(vals, cty.ObjectVal(obj))

	resp.Result = cty.ObjectVal(map[string]cty.Value{
		"data":   cty.TupleVal(vals),
		"config": req.Config,
	})
	return
}

func (s simple) ValidateStateStoreConfig(req providers.ValidateStateStoreConfigRequest) providers.ValidateStateStoreConfigResponse {
	panic("not implemented")
}

func (s simple) ConfigureStateStore(req providers.ConfigureStateStoreRequest) providers.ConfigureStateStoreResponse {
	panic("not implemented")
}

func (s simple) GetStates(req providers.GetStatesRequest) providers.GetStatesResponse {
	panic("not implemented")
}

func (s simple) DeleteState(req providers.DeleteStateRequest) providers.DeleteStateResponse {
	panic("not implemented")
}

func (s simple) PlanAction(providers.PlanActionRequest) providers.PlanActionResponse {
	// Our schema doesn't include any actions, so it should be
	// impossible to get here.
	panic("PlanAction on provider that didn't declare any actions")
}

func (s simple) InvokeAction(providers.InvokeActionRequest) providers.InvokeActionResponse {
	// Our schema doesn't include any actions, so it should be
	// impossible to get here.
	panic("InvokeAction on provider that didn't declare any actions")
}

func (s simple) ValidateActionConfig(providers.ValidateActionConfigRequest) providers.ValidateActionConfigResponse {
	// Our schema doesn't include any actions, so it should be
	// impossible to get here.
	panic("ValidateActionConfig on provider that didn't declare any actions")
}

func (s simple) Close() error {
	return nil
}
