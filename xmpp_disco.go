package xmpp

import (
	"encoding/xml"
)

const (
	XMPPNS_DISCO_ITEMS = "http://jabber.org/protocol/disco#items"
	XMPPNS_DISCO_INFO  = "http://jabber.org/protocol/disco#info"
)

type clientDiscoFeature struct {
	XMLName xml.Name `xml:"feature"`
	Var     string   `xml:"var,attr"`
}

type clientDiscoIdentity struct {
	XMLName  xml.Name `xml:"identity"`
	Category string   `xml:"category,attr"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
}

type clientDiscoQuery struct {
	XMLName    xml.Name              `xml:"query"`
	Features   []clientDiscoFeature  `xml:"feature"`
	Identities []clientDiscoIdentity `xml:"identity"`
}

type DiscoIdentity struct {
	Category string
	Type     string
	Name     string
}

type DiscoResult struct {
	Features   []string
	Identities []DiscoIdentity
}

func clientFeaturesToReturn(features []clientDiscoFeature) []string {
	var ret []string

	for _, feature := range features {
		ret = append(ret, feature.Var)
	}

	return ret
}

func clientIdentitiesToReturn(identities []clientDiscoIdentity) []DiscoIdentity {
	var ret []DiscoIdentity

	for _, id := range identities {
		ret = append(ret, DiscoIdentity{
			Category: id.Category,
			Type:     id.Type,
			Name:     id.Name,
		})
	}

	return ret
}
