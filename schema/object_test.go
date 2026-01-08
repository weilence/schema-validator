package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
)

func TestObjectSchema_NilPointerFields(t *testing.T) {
	type SNMPInfo struct {
		Protocol  string
		Version   string
		Community string
	}

	type SyslogInfo struct {
		Protocol string
		Encoding string
	}

	type KafkaInfo struct {
		Topic    string
		Security string
	}

	type CommonInfo struct {
		SNMP   *SNMPInfo
		Syslog *SyslogInfo
		Kafka  *KafkaInfo
	}

	type Request struct {
		DeviceName string
		NotifyType string
		Info       CommonInfo
	}

	tests := []struct {
		name    string
		data    Request
		wantErr bool
	}{
		{
			name: "all fields nil",
			data: Request{
				DeviceName: "device1",
				NotifyType: "SNMP",
				Info: CommonInfo{
					SNMP:   nil,
					Syslog: nil,
					Kafka:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "only SNMP set, others nil",
			data: Request{
				DeviceName: "device1",
				NotifyType: "SNMP",
				Info: CommonInfo{
					SNMP: &SNMPInfo{
						Protocol:  "UDP",
						Version:   "v2c",
						Community: "public",
					},
					Syslog: nil,
					Kafka:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "only Syslog set, others nil",
			data: Request{
				DeviceName: "device1",
				NotifyType: "Syslog",
				Info: CommonInfo{
					SNMP: nil,
					Syslog: &SyslogInfo{
						Protocol: "TCP",
						Encoding: "UTF-8",
					},
					Kafka: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "only Kafka set, others nil",
			data: Request{
				DeviceName: "device1",
				NotifyType: "Kafka",
				Info: CommonInfo{
					SNMP:   nil,
					Syslog: nil,
					Kafka: &KafkaInfo{
						Topic:    "alarms",
						Security: "SASL/PLAIN",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "all fields set",
			data: Request{
				DeviceName: "device1",
				NotifyType: "SNMP",
				Info: CommonInfo{
					SNMP: &SNMPInfo{
						Protocol:  "UDP",
						Version:   "v3",
						Community: "private",
					},
					Syslog: &SyslogInfo{
						Protocol: "UDP",
						Encoding: "UTF-8",
					},
					Kafka: &KafkaInfo{
						Topic:    "logs",
						Security: "SASL/SCRAM",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objSchema := NewObject()
			objSchema.AddField("DeviceName", NewField())
			objSchema.AddField("NotifyType", NewField())

			infoSchema := NewObject()
			infoSchema.AddField("SNMP", NewObject())
			infoSchema.AddField("Syslog", NewObject())
			infoSchema.AddField("Kafka", NewObject())
			objSchema.AddField("Info", infoSchema)

			accessor := data.New(tt.data)
			ctx := NewContext(objSchema, accessor)

			err := objSchema.Validate(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectSchema_NilPointerValidation(t *testing.T) {
	type Inner struct {
		Value string
	}

	type Outer struct {
		Field *Inner
	}

	t.Run("nil pointer should not cause error", func(t *testing.T) {
		testData := Outer{Field: nil}

		objSchema := NewObject()
		innerSchema := NewObject()
		innerSchema.AddField("Value", NewField())
		objSchema.AddField("Field", innerSchema)

		accessor := data.New(testData)
		ctx := NewContext(objSchema, accessor)

		err := objSchema.Validate(ctx)
		assert.NoError(t, err)
	})

	t.Run("non-nil pointer should validate nested fields", func(t *testing.T) {
		testData := Outer{
			Field: &Inner{Value: "test"},
		}

		objSchema := NewObject()
		innerSchema := NewObject()
		innerSchema.AddField("Value", NewField())
		objSchema.AddField("Field", innerSchema)

		accessor := data.New(testData)
		ctx := NewContext(objSchema, accessor)

		err := objSchema.Validate(ctx)
		assert.NoError(t, err)
	})
}

func TestObjectSchema_ModifySchemaWithNilPointers(t *testing.T) {
	type SNMPInfo struct {
		Protocol  string
		Version   string
		Community string
	}

	type SyslogInfo struct {
		Protocol string
	}

	type CommonInfo struct {
		SNMP   *SNMPInfo
		Syslog *SyslogInfo
	}

	type Request struct {
		NotifyType string
		Info       CommonInfo
	}

	testData := Request{
		NotifyType: "SNMP",
		Info: CommonInfo{
			SNMP: &SNMPInfo{
				Protocol:  "UDP",
				Version:   "v2c",
				Community: "public",
			},
			Syslog: nil,
		},
	}

	objSchema := NewObject()
	objSchema.AddField("NotifyType", NewField())

	infoSchema := NewObject()
	infoSchema.AddField("SNMP", NewObject())
	infoSchema.AddField("Syslog", NewObject())
	objSchema.AddField("Info", infoSchema)

	infoSchema.RemoveField("Syslog")

	accessor := data.New(testData)
	ctx := NewContext(objSchema, accessor)

	err := objSchema.Validate(ctx)
	assert.NoError(t, err)
}

func TestObjectSchema_RealWorldAlarmNotifyCase(t *testing.T) {
	type SNMPInfo struct {
		Protocol  string `json:"protocol"`
		Version   string `json:"version"`
		Community string `json:"community"`
	}

	type SyslogInfo struct {
		Protocol string `json:"protocol"`
		Encoding string `json:"encoding"`
	}

	type KafkaInfo struct {
		Topic    string `json:"topic"`
		Security string `json:"security"`
	}

	type CommonInfo struct {
		SNMP   *SNMPInfo   `json:"snmp"`
		Syslog *SyslogInfo `json:"syslog"`
		Kafka  *KafkaInfo  `json:"kafka"`
	}

	type AlarmNotifyRequest struct {
		DeviceName      string     `json:"deviceName"`
		NotifyType      string     `json:"notifyType"`
		NotifyAddresses []string   `json:"notifyAddresses"`
		AlarmTypes      []string   `json:"alarmTypes"`
		AlarmLevels     []string   `json:"alarmLevels"`
		Enabled         bool       `json:"enabled"`
		Info            CommonInfo `json:"info"`
	}

	t.Run("SNMP notification with nil Syslog and Kafka", func(t *testing.T) {
		req := AlarmNotifyRequest{
			DeviceName:      "global",
			NotifyType:      "SNMP",
			NotifyAddresses: []string{"1.1.1.1"},
			AlarmTypes:      []string{"CPU Temperature", "CPU Usage"},
			AlarmLevels:     []string{"critical", "error", "warn", "info"},
			Enabled:         true,
			Info: CommonInfo{
				SNMP: &SNMPInfo{
					Protocol:  "UDP",
					Version:   "v2c",
					Community: "public",
				},
				Syslog: nil,
				Kafka:  nil,
			},
		}

		objSchema := NewObject()
		objSchema.AddField("DeviceName", NewField())
		objSchema.AddField("NotifyType", NewField())
		objSchema.AddField("NotifyAddresses", NewArray(NewField()))
		objSchema.AddField("AlarmTypes", NewArray(NewField()))
		objSchema.AddField("AlarmLevels", NewArray(NewField()))
		objSchema.AddField("Enabled", NewField())

		infoSchema := NewObject()
		snmpSchema := NewObject()
		snmpSchema.AddField("Protocol", NewField())
		snmpSchema.AddField("Version", NewField())
		snmpSchema.AddField("Community", NewField())
		infoSchema.AddField("SNMP", snmpSchema)
		infoSchema.AddField("Syslog", NewObject())
		infoSchema.AddField("Kafka", NewObject())
		objSchema.AddField("Info", infoSchema)

		accessor := data.New(req)
		ctx := NewContext(objSchema, accessor)

		err := objSchema.Validate(ctx)
		assert.NoError(t, err, "Should handle nil pointer fields in CommonInfo struct")
	})

	t.Run("Syslog notification with nil SNMP and Kafka", func(t *testing.T) {
		req := AlarmNotifyRequest{
			DeviceName:      "device-1",
			NotifyType:      "Syslog",
			NotifyAddresses: []string{"192.168.1.100:514"},
			AlarmTypes:      []string{"Memory Usage"},
			AlarmLevels:     []string{"error", "critical"},
			Enabled:         true,
			Info: CommonInfo{
				SNMP: nil,
				Syslog: &SyslogInfo{
					Protocol: "TCP",
					Encoding: "UTF-8",
				},
				Kafka: nil,
			},
		}

		objSchema := NewObject()
		objSchema.AddField("DeviceName", NewField())
		objSchema.AddField("NotifyType", NewField())
		objSchema.AddField("NotifyAddresses", NewArray(NewField()))
		objSchema.AddField("AlarmTypes", NewArray(NewField()))
		objSchema.AddField("AlarmLevels", NewArray(NewField()))
		objSchema.AddField("Enabled", NewField())

		infoSchema := NewObject()
		infoSchema.AddField("SNMP", NewObject())
		syslogSchema := NewObject()
		syslogSchema.AddField("Protocol", NewField())
		syslogSchema.AddField("Encoding", NewField())
		infoSchema.AddField("Syslog", syslogSchema)
		infoSchema.AddField("Kafka", NewObject())
		objSchema.AddField("Info", infoSchema)

		accessor := data.New(req)
		ctx := NewContext(objSchema, accessor)

		err := objSchema.Validate(ctx)
		assert.NoError(t, err, "Should handle nil pointer fields in CommonInfo struct")
	})

	t.Run("all notification types nil", func(t *testing.T) {
		req := AlarmNotifyRequest{
			DeviceName:      "device-2",
			NotifyType:      "Email",
			NotifyAddresses: []string{"admin@example.com"},
			AlarmTypes:      []string{"Disk Usage"},
			AlarmLevels:     []string{"warn"},
			Enabled:         false,
			Info: CommonInfo{
				SNMP:   nil,
				Syslog: nil,
				Kafka:  nil,
			},
		}

		objSchema := NewObject()
		objSchema.AddField("DeviceName", NewField())
		objSchema.AddField("NotifyType", NewField())
		objSchema.AddField("NotifyAddresses", NewArray(NewField()))
		objSchema.AddField("AlarmTypes", NewArray(NewField()))
		objSchema.AddField("AlarmLevels", NewArray(NewField()))
		objSchema.AddField("Enabled", NewField())

		infoSchema := NewObject()
		infoSchema.AddField("SNMP", NewObject())
		infoSchema.AddField("Syslog", NewObject())
		infoSchema.AddField("Kafka", NewObject())
		objSchema.AddField("Info", infoSchema)

		accessor := data.New(req)
		ctx := NewContext(objSchema, accessor)

		err := objSchema.Validate(ctx)
		assert.NoError(t, err, "Should handle all nil pointer fields in CommonInfo struct")
	})
}
