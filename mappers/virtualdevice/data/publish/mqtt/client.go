package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubeedge/mapper-generator/pkg/common"
	"github.com/kubeedge/mapper-generator/pkg/global"
	"k8s.io/klog/v2"
	"os"
)

type PushMethod struct {
	MQTT *MQTTConfig `json:"http"`
}

type MQTTConfig struct {
	Address  string `json:"address,omitempty"`
	Topic    string `json:"topic,omitempty"`
	QoS      int    `json:"qos,omitempty"`
	Retained bool   `json:"retained,omitempty"`
}

func NewDataPanel(config json.RawMessage) (global.DataPanel, error) {
	mqttConfig := new(MQTTConfig)
	err := json.Unmarshal(config, mqttConfig)
	if err != nil {
		return nil, err
	}
	return &PushMethod{
		MQTT: mqttConfig,
	}, nil
}

func (pm *PushMethod) InitPushMethod() error {
	// TODO add init code
	fmt.Println("Init Mqtt")
	return nil
}

func (pm *PushMethod) Push(data *common.DataModel) {
	// TODO add push code
	klog.V(1).Infof("Publish %v to %s on topic: %s, Qos: %d, Retained: %v",
		data.Value, pm.MQTT.Address, pm.MQTT.Topic, pm.MQTT.QoS, pm.MQTT.Retained)

	opts := mqtt.NewClientOptions().AddBroker(pm.MQTT.Address)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	token := client.Publish(pm.MQTT.Topic, byte(pm.MQTT.QoS), pm.MQTT.Retained, data.Value)
	token.Wait()

	client.Disconnect(250)
	klog.V(1).Info("###############")
	klog.V(1).Info("Message published.")
}