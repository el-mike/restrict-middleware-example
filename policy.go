package main

import "github.com/el-mike/restrict/v2"

var Policy = &restrict.PolicyDefinition{
	PermissionPresets: restrict.PermissionPresets{
		"updateOwn": &restrict.Permission{
			Action: "update",
			Conditions: restrict.Conditions{
				&restrict.EqualCondition{
					ID: "isOwner",
					Left: &restrict.ValueDescriptor{
						Source: restrict.ResourceField,
						Field:  "CreatedBy",
					},
					Right: &restrict.ValueDescriptor{
						Source: restrict.SubjectField,
						Field:  "ID",
					},
				},
			},
		},
	},
	Roles: restrict.Roles{
		"User": {
			Description: "This is a simple User role, with permissions for basic chat operations.",
			Grants: restrict.GrantsMap{
				"Conversation": {
					&restrict.Permission{Action: "read"},
					&restrict.Permission{Action: "create"},
					&restrict.Permission{Preset: "updateOwn"},
					&restrict.Permission{
						Action: "delete",
						Conditions: restrict.Conditions{
							&restrict.EmptyCondition{
								ID: "deleteActive",
								Value: &restrict.ValueDescriptor{
									Source: restrict.ResourceField,
									Field:  "Active",
								},
							},
						},
					},
				},
			},
		},
		"Admin": {
			Description: "This is an Admin role, with permissions to manage Users.",
			Parents:     []string{"User"},
			Grants: restrict.GrantsMap{
				"User": {
					&restrict.Permission{Action: "read"},
					&restrict.Permission{Action: "create"},
				},
			},
		},
	},
}
