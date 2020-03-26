package xmpp

import (
	"fmt"
	"strconv"
	"encoding/xml"
)

const (
	XMPPNS_DISCO_ITEMS = "http://jabber.org/protocol/disco#items"
	XMPPNS_DISCO_INFO = "http://jabber.org/protocol/disco#info"

	XMPPNS_DISCO_COMMANDS = "http://jabber.org/protocol/commands"

	IQTypeGet = "get"
	IQTypeSet = "set"
	IQTypeResult = "result"
)

type clientDiscoFeature struct {
	XMLName xml.Name `xml:"feature"`
	Var string `xml:"var,attr"`
}

type clientDiscoIdentity struct {
	XMLName xml.Name `xml:"identity"`
	Category string `xml:"category,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
}

type clientDiscoQuery struct {
	XMLName xml.Name `xml:"query"`
	Features []clientDiscoFeature `xml:"feature"`
	Identities []clientDiscoIdentity `xml:"identity"`
}

type clientDiscoQueryItem struct {
	XMLName xml.Name `xml:"item"`

	JID string `xml:"jid,attr"`
	Name string `xml:"name,attr"`
	Node string `xml:"node,attr"`
}

type clientDiscoQueryItems struct {
	XMLName xml.Name `xml:"query"`

	Items []clientDiscoQueryItem `xml:"item"`
}

type DiscoIdentity struct {
	Category string
	Type string
	Name string
}

type DiscoResult struct {
	Features   []string
	Identities []DiscoIdentity
}

func (c *Client) Discovery() (string, error) {
	// use getCookie for a pseudo random id.
	reqID := strconv.FormatUint(uint64(getCookie()), 10)
	return c.RawInformationQuery(c.jid, c.domain, reqID, IQTypeGet, XMPPNS_DISCO_ITEMS, "")
}

// Discover the capabilities of a node from the server according to XEP-0030
func (c *Client) DiscoverNode(node string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s' node='%s'/>", XMPPNS_DISCO_INFO, node)
	return c.RawInformation(c.jid, c.domain, "info3", IQTypeGet, query)
}

// Discover the capabilities of the server according to XEP-0030
func (c *Client) DiscoverNodeItems(node string) (string, error) {
	query := fmt.Sprintf("<query xmlns='%s' node='%s'/>", XMPPNS_DISCO_ITEMS, node)
	return c.RawInformation(c.jid, c.domain, "items1", IQTypeGet, query)
}

func (c *Client) DiscoverItems() (string, error) {
	query := fmt.Sprintf("<query xmlns='%s'/>", XMPPNS_DISCO_INFO)
	return c.RawInformation(c.jid, c.domain, "info1", IQTypeGet, query)
}

// RawInformationQuery sends an information query request to the server.
func (c *Client) RawInformationQuery(from, to, id, iqType, requestNamespace, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'><query xmlns='%s'>%s</query></iq>"
	_, err := fmt.Fprintf(c.conn, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, requestNamespace, body)
	return id, err
}

// rawInformation send a IQ request with the the payload body to the server
func (c *Client) RawInformation(from, to, id, iqType, body string) (string, error) {
	const xmlIQ = "<iq from='%s' to='%s' id='%s' type='%s'>%s</iq>"
	_, err := fmt.Fprintf(c.conn, xmlIQ, xmlEscape(from), xmlEscape(to), id, iqType, body)
	return id, err
}
