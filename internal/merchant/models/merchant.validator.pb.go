// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: internal/merchant/merchant.proto

package models

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/infobloxopen/protoc-gen-gorm/options"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *MerchantAccount) Validate() error {
	for _, item := range this.Owners {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Owners", err)
			}
		}
	}
	if this.Address != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Address); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Address", err)
		}
	}
	for _, item := range this.ItemsOrServicesSold {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ItemsOrServicesSold", err)
			}
		}
	}
	if this.ShopSettings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ShopSettings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ShopSettings", err)
		}
	}
	for _, item := range this.Tags {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Tags", err)
			}
		}
	}
	return nil
}
func (this *Settings) Validate() error {
	if this.PaymentDetails != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PaymentDetails); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PaymentDetails", err)
		}
	}
	for _, item := range this.ShopPolicy {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ShopPolicy", err)
			}
		}
	}
	for _, item := range this.PrivacyPolicy {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PrivacyPolicy", err)
			}
		}
	}
	if this.ReturnPolicy != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ReturnPolicy); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ReturnPolicy", err)
		}
	}
	if this.ShippingPolicy != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.ShippingPolicy); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("ShippingPolicy", err)
		}
	}
	return nil
}
func (this *Settings_PaymentDetails) Validate() error {
	return nil
}
func (this *ItemSold) Validate() error {
	return nil
}
func (this *Address) Validate() error {
	return nil
}
func (this *Owner) Validate() error {
	return nil
}
func (this *Tags) Validate() error {
	return nil
}
func (this *Policy) Validate() error {
	for _, item := range this.Tags {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Tags", err)
			}
		}
	}
	return nil
}
func (this *ReturnPolicy) Validate() error {
	if this.PolicyMeta != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PolicyMeta); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PolicyMeta", err)
		}
	}
	return nil
}
func (this *ShippingPolicy) Validate() error {
	if this.PolicyMeta != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PolicyMeta); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PolicyMeta", err)
		}
	}
	return nil
}
