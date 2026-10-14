package astgen_test

import (
	"testing"

	"github.com/scjalliance/astconf/astgen"
	"github.com/scjalliance/astconf/astorg"
)

// legacyPagingData returns the multi-server data set with more paging
// groups. Paging groups reach phones through people, roles and tags as well
// as locations. A softphone belongs to both a person and a role, and a
// phone at a location is reassigned to a role. One group has no members and
// one has no alert.
func legacyPagingData() *astorg.DataSet {
	data := multiServerData()
	data.Tags = astorg.TagList{{Name: "ops", PagingGroups: []string{"510", "500"}}, {Name: "none"}}
	data.PagingGroups = append(data.PagingGroups,
		astorg.PagingGroup{Name: "Ops", Extension: "510"},
		astorg.PagingGroup{Name: "Empty", Extension: "520", Alert: "x"},
		astorg.PagingGroup{Name: "Yard", Extension: "530", Alert: "yard"},
	)
	data.People[0].PagingGroups = []string{"530", "510"}
	data.People[0].Tags = []string{"ops", "missing"}
	data.People[2].Tags = []string{"ops"}
	data.PhoneRoles[0].PagingGroups = []string{"530"}
	data.PhoneRoles[1].PagingGroups = []string{"510", "530"}
	data.PhoneRoles[1].Tags = []string{"ops"}
	data.PhoneRoles[1].Softphones = append(data.PhoneRoles[1].Softphones, "fred.soft", "wilma.soft")
	data.PhoneRoles[1].Phones = []string{"000fd3000003", "000fd3000001"}
	data.Locations[2].PagingGroups = []string{"530"}
	data.Phones = append(data.Phones, astorg.Phone{MAC: "000fd3000009", Location: "Remote"})
	return data
}

// TestPagingGroupExtensionsLegacy pins the output of PagingGroupExtensions.
// The expected strings were rendered by the generator before its
// membership logic moved into pagingMembers, so the shared helper cannot
// change what the old generator renders.
func TestPagingGroupExtensionsLegacy(t *testing.T) {
	tests := []struct {
		name string
		data *astorg.DataSet
		want string
	}{
		{
			name: "no paging groups",
			data: testData(),
			want: "[paging]\n",
		},
		{
			name: "locations only",
			data: multiServerData(),
			want: "[paging]\n" +
				"\n" +
				"; All Call\n" +
				"exten => 500,1,Noop(Paging All Call)\n" +
				"same => n,Playback(beep)\n" +
				"same => n,SIPAddHeader(Alert-Info:intercom)\n" +
				"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
				"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
				"same => n,Page(SIP/000fd3000003&SIP/barney&SIP/fred&SIP/lobby,is)\n" +
				"same => n,Hangup()\n",
		},
		{
			name: "people, roles and tags",
			data: legacyPagingData(),
			want: "[paging]\n" +
				"\n" +
				"; All Call\n" +
				"exten => 500,1,Noop(Paging All Call)\n" +
				"same => n,Playback(beep)\n" +
				"same => n,SIPAddHeader(Alert-Info:intercom)\n" +
				"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
				"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
				"same => n,Page(SIP/barney&SIP/fred&SIP/fred.soft&SIP/lobby&SIP/wilma.soft&SIP/yard&SIP/yard.soft,is)\n" +
				"same => n,Hangup()\n" +
				"\n" +
				"; Ops\n" +
				"exten => 510,1,Noop(Paging Ops)\n" +
				"same => n,Playback(beep)\n" +
				"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
				"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
				"same => n,Page(SIP/fred&SIP/fred.soft&SIP/wilma.soft&SIP/yard&SIP/yard.soft,is)\n" +
				"same => n,Hangup()\n" +
				"\n" +
				"; Empty\n" +
				"exten => 520,1,Noop(Paging Empty)\n" +
				"same => n,Playback(beep)\n" +
				"same => n,SIPAddHeader(Alert-Info:x)\n" +
				"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
				"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
				"same => n,Hangup()\n" +
				"\n" +
				"; Yard\n" +
				"exten => 530,1,Noop(Paging Yard)\n" +
				"same => n,Playback(beep)\n" +
				"same => n,SIPAddHeader(Alert-Info:yard)\n" +
				"same => n,Set(CALLERID(name)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(name)}:Intercom)})\n" +
				"same => n,Set(CALLERID(num)=${IF($[${LEN(${CALLERID(num)})}>0]?${CALLERID(num)}:Intercom)})\n" +
				"same => n,Page(SIP/000fd3000009&SIP/fred&SIP/fred.soft&SIP/lobby&SIP/yard&SIP/yard.soft,is)\n" +
				"same => n,Hangup()\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := astgen.PagingGroupExtensions(tt.data, "paging")
			if got := render(t, &section); got != tt.want {
				t.Errorf("PagingGroupExtensions() rendered:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}
