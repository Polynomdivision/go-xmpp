package xmpp

import (
	"encoding/xml"
	"fmt"
)

type AdhocCommand struct {
	JID string
	Name string
	Node string
}

type AdhocResult struct {
	Node string
	Status string

	Data []byte
}

type clientUptimeNote struct {
	XMLName xml.Name `xml:"note"`

	InnerXML []byte `xml:",innerxml"`
}

type Uptime struct {
	Uptime string
}

type clientAdhocCommand struct {
	XMLName xml.Name `xml:"command"`
	SessionID string `xml:"sessionid,attr"`
	Node string `xml:"node,attr"`
	Status string `xml:"status,attr"`

	//X clientAdhocX `xml:"x"`
	Data []byte `xml:",innerxml"`
}

// type clientAdhocX struct {
// 	XMLName xml.Name `xml:"x"`
// 	XMLNS string `xml:"xmlns,attr"`
// 	Result string `xml:"result,attr"`

// 	InnerXML []byte `xml:",innerxml"`
// }

func (c *Client) AdhocGetCommands() {
	c.DiscoverNodeItems(XMPPNS_DISCO_COMMANDS)
}

func (c *Client) AdhocExecuteCommand(node string) {
	c.RawInformation(c.jid, c.domain, "exec1", "set", adhocCommandStanza(node, "execute"))
}

func adhocCommandStanza(node, action string) string {
	return fmt.Sprintf("<command xmlns='%s' node='%s' action='%s'/>",
		XMPPNS_DISCO_COMMANDS, node, action)
}

func adhocIsCommandList(query XMLElement) bool {
	return query.XMLNS == XMPPNS_DISCO_ITEMS && query.Node == XMPPNS_DISCO_COMMANDS
}

func adhocParseCommandList(innerXML string) ([]AdhocCommand, error) {
	var disco clientDiscoQueryItems
	err := xml.Unmarshal([]byte("<query>" + innerXML + "</query>"), &disco)
	if err != nil {
		return []AdhocCommand{}, err
	}
	
	var cmds []AdhocCommand
	for _, i := range disco.Items {
		cmds = append(cmds, AdhocCommand{
			i.JID,
			i.Name,
			i.Node,
		})
	}

	return cmds, nil
}
