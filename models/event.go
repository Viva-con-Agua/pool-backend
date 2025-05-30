package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	//EventCreate represents the model for creating an event.
	EventCreate struct {
		Name                  string           `json:"name" bson:"name"`
		TypeOfEvent           string           `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string           `json:"additional_information" bson:"additional_information"`
		Website               string           `json:"website" bson:"website"`
		TourID                string           `json:"tour_id" bson:"tour_id"`
		Location              Location         `json:"location" bson:"location"`
		MeetingURL            string           `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string         `json:"artist_ids" bson:"artist_ids"`
		OrganizerID           string           `json:"organizer_id" bson:"organizer_id"`
		StartAt               int64            `json:"start_at" bson:"start_at"`
		EndAt                 int64            `json:"end_at" bson:"end_at"`
		CrewID                string           `json:"crew_id" bson:"crew_id"`
		OrganisationID        string           `json:"organisation_id" bson:"organisation_id"`
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
	}
	EventDatabase struct {
		ID                    string            `json:"id" bson:"_id"`
		Name                  string            `json:"name" bson:"name"`
		TypeOfEvent           string            `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string            `json:"additional_information" bson:"additional_information"`
		Website               string            `json:"website" bson:"website"`
		TourID                string            `json:"tour_id" bson:"tour_id"`
		Location              Location          `json:"location" bson:"location"`
		MeetingURL            string            `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string          `json:"artist_ids" bson:"artist_ids"`
		OrganizerID           string            `json:"organizer_id" bson:"organizer_id"`
		StartAt               int64             `json:"start_at" bson:"start_at"`
		EndAt                 int64             `json:"end_at" bson:"end_at"`
		CrewID                string            `json:"crew_id" bson:"crew_id"`
		OrganisationID        string            `json:"organisation_id" bson:"organisation_id"`
		TakingID              string            `json:"taking_id" bson:"taking_id"`
		EventASPID            string            `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string            `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal      `json:"external_asp" bson:"external_asp"`
		Application           EventApplication  `json:"application" bson:"application"`
		Applications          EventApplications `json:"applications" bson:"applications"`
		EventTools            EventTools        `json:"event_tools" bson:"event_tools"`
		CreatorID             string            `json:"creator_id" bson:"creator_id"`
		EventState            EventState        `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified     `json:"modified" bson:"modified"`
	}
	EventTools struct {
		Tools   []string `json:"tools" bson:"tools"`
		Special string   `json:"special" bson:"special"`
	}
	//EventApplication represents the application type of an event.
	EventApplication struct {
		StartDate      int64 `json:"start_date" bson:"start_date"`
		EndDate        int64 `json:"end_date" bson:"end_date"`
		SupporterCount int   `json:"supporter_count" bson:"supporter_count"`
	}
	EventApplications struct {
		Total     int `json:"total" bson:"total"`
		Confirmed int `json:"confirmed" bson:"confirmed"`
		Rejected  int `json:"rejected" bson:"rejected"`
		Requested int `json:"requested" bson:"requested"`
		Withdrawn int `json:"withdrawn" bson:"withdrawn"`
	}
	//EventState represents the state of an event.
	EventState struct {
		State                string `json:"state" bson:"state"`
		CrewConfirmation     string `json:"crew_confirmation" bson:"crew_confirmation"`
		InternalConfirmation string `json:"internal_confirmation" bson:"internal_confirmation"`
		TakingID             string `json:"taking_id" bson:"taking_id"`
		DepositID            string `json:"deposit_id" bson:"deposit_id"`
		OldState             string `json:"-" bson:"old_state"`
	}
	EventStatePublic struct {
		State string `json:"state" bson:"state"`
	}
	EventDetails struct {
		MeetingURL string         `json:"meeting_url" bson:"meeting_url"`
		EventASP   EventASPPublic `json:"event_asp" bson:"event_asp"`
	}
	EventASPPublic struct {
		FullName    string `bson:"full_name" json:"full_name"`
		DisplayName string `bson:"display_name" json:"display_name"`
		Phone       string `bson:"phone" json:"phone"`
		Email       string `json:"email" bson:"email"`
		Mattermost  string `bson:"mattermost_username" json:"mattermost_username"`
	}
	EventPublic struct {
		ID                    string           `json:"id" bson:"_id"`
		Name                  string           `json:"name" bson:"name"`
		TypeOfEvent           string           `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string           `json:"additional_information" bson:"additional_information"`
		Website               string           `json:"website" bson:"website"`
		TourID                string           `json:"tour_id" bson:"tour_id"`
		Location              Location         `json:"location" bson:"location"`
		ArtistIDs             []string         `json:"artist_ids" bson:"artist_ids"`
		Artists               []Artist         `json:"artists" bson:"artists"`
		OrganizerID           string           `json:"organizer_id" bson:"organizer_id"`
		Organizer             Organizer        `json:"organizer" bson:"organizer"`
		Organisation          Organisation     `json:"organisation" bson:"organisation"`
		StartAt               int64            `json:"start_at" bson:"start_at"`
		EndAt                 int64            `json:"end_at" bson:"end_at"`
		CrewID                string           `json:"crew_id" bson:"crew_id"`
		Crew                  CrewPublic       `json:"crew" bson:"crew"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		EventState            EventStatePublic `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified    `json:"modified" bson:"modified"`
	}
	ListEvent struct {
		ID                    string                 `json:"id" bson:"_id"`
		Name                  string                 `json:"name" bson:"name"`
		TypeOfEvent           string                 `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string                 `json:"additional_information" bson:"additional_information"`
		Website               string                 `json:"website" bson:"website"`
		TourID                string                 `json:"tour_id" bson:"tour_id"`
		Location              Location               `json:"location" bson:"location"`
		ArtistIDs             []string               `json:"artist_ids" bson:"artist_ids"`
		Artists               []Artist               `json:"artists" bson:"artists"`
		OrganizerID           string                 `json:"organizer_id" bson:"organizer_id"`
		Organizer             Organizer              `json:"organizer" bson:"organizer"`
		Organisation          Organisation           `json:"organisation" bson:"organisation"`
		StartAt               int64                  `json:"start_at" bson:"start_at"`
		EndAt                 int64                  `json:"end_at" bson:"end_at"`
		CrewID                string                 `json:"crew_id" bson:"crew_id"`
		Crew                  CrewPublic             `json:"crew" bson:"crew"`
		EventASPID            string                 `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string                 `json:"internal_asp_id" bson:"internal_asp_id"`
		Application           EventApplication       `json:"application" bson:"application"`
		Applications          EventApplications      `json:"applications" bson:"applications"`
		Participation         []ParticipationMinimal `json:"participations" bson:"participations"`
		EventTools            EventTools             `json:"event_tools" bson:"event_tools"`
		EventState            EventState             `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified          `json:"modified" bson:"modified"`
	}
	ListDetailsEvent struct {
		ID                    string                 `json:"id" bson:"_id"`
		Name                  string                 `json:"name" bson:"name"`
		TypeOfEvent           string                 `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string                 `json:"additional_information" bson:"additional_information"`
		Website               string                 `json:"website" bson:"website"`
		TourID                string                 `json:"tour_id" bson:"tour_id"`
		Location              Location               `json:"location" bson:"location"`
		MeetingURL            string                 `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string               `json:"artist_ids" bson:"artist_ids"`
		Artists               []Artist               `json:"artists" bson:"artists"`
		OrganizerID           string                 `json:"organizer_id" bson:"organizer_id"`
		Organizer             Organizer              `json:"organizer" bson:"organizer"`
		StartAt               int64                  `json:"start_at" bson:"start_at"`
		EndAt                 int64                  `json:"end_at" bson:"end_at"`
		CrewID                string                 `json:"crew_id" bson:"crew_id"`
		Crew                  Crew                   `json:"crew" bson:"crew"`
		EventASPID            string                 `json:"event_asp_id" bson:"event_asp_id"`
		EventASP              EventASPPublic         `json:"event_asp" bson:"event_asp"`
		InternalASPID         string                 `json:"internal_asp_id" bson:"internal_asp_id"`
		Application           EventApplication       `json:"application" bson:"application"`
		Applications          EventApplications      `json:"applications" bson:"applications"`
		Participation         []ParticipationMinimal `json:"participations" bson:"participations"`
		EventTools            EventTools             `json:"event_tools" bson:"event_tools"`
		EventState            EventState             `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified          `json:"modified" bson:"modified"`
	}
	Event struct {
		ID                    string            `json:"id" bson:"_id"`
		Name                  string            `json:"name" bson:"name"`
		TypeOfEvent           string            `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string            `json:"additional_information" bson:"additional_information"`
		Website               string            `json:"website" bson:"website"`
		TourID                string            `json:"tour_id" bson:"tour_id"`
		Location              Location          `json:"location" bson:"location"`
		MeetingURL            string            `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string          `json:"artist_ids" bson:"artist_ids"`
		Artists               []Artist          `json:"artists" bson:"artists"`
		OrganizerID           string            `json:"organizer_id" bson:"organizer_id"`
		Organizer             Organizer         `json:"organizer" bson:"organizer"`
		StartAt               int64             `json:"start_at" bson:"start_at"`
		EndAt                 int64             `json:"end_at" bson:"end_at"`
		CrewID                string            `json:"crew_id" bson:"crew_id"`
		Crew                  Crew              `json:"crew" bson:"crew"`
		EventASPID            string            `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string            `json:"internal_asp_id" bson:"internal_asp_id"`
		EventASP              User              `json:"event_asp" bson:"event_asp"`
		InteralASP            User              `json:"internal_asp" bson:"internal_asp"`
		OrganisationID        string            `json:"organisation_id" bson:"organisation_id"`
		ExternalASP           UserExternal      `json:"external_asp" bson:"external_asp"`
		TakingID              string            `json:"taking_id" bson:"taking_id"`
		DepositID             string            `json:"deposit_id" bson:"deposit_id"`
		Application           EventApplication  `json:"application" bson:"application"`
		Applications          EventApplications `json:"applications" bson:"applications"`
		Participation         []Participation   `json:"participations" bson:"participations"`
		EventTools            EventTools        `json:"event_tools" bson:"event_tools"`
		Creator               User              `json:"creator" bson:"creator"`
		EventState            EventState        `json:"event_state" bson:"event_state"`
		EditorID              string            `json:"editor_id" bson:"editor_id"`
		Modified              vmod.Modified     `json:"modified" bson:"modified"`
	}
	EventUpdate struct {
		ID                    string           `json:"id" bson:"_id"`
		Name                  string           `json:"name" bson:"name"`
		TypeOfEvent           string           `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string           `json:"additional_information" bson:"additional_information"`
		Website               string           `json:"website" bson:"website"`
		TourID                string           `json:"tour_id" bson:"tour_id"`
		Location              Location         `json:"location" bson:"location"`
		MeetingURL            string           `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string         `json:"artist_ids" bson:"artist_ids"`
		OrganizerID           string           `json:"organizer_id" bson:"organizer_id"`
		StartAt               int64            `json:"start_at" bson:"start_at"`
		EndAt                 int64            `json:"end_at" bson:"end_at"`
		CrewID                string           `json:"crew_id" bson:"crew_id"`
		OrganisationID        string           `json:"organisation_id" bson:"organisation_id"`
		Crew                  Crew             `json:"crew" bson:"crew"`
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
	}
	EventApplicationsUpdate struct {
		ID           string            `json:"id" bson:"_id"`
		Applications EventApplications `json:"applications" bson:"applications"`
	}
	EventParam struct {
		ID string `param:"id"`
	}

	EventQuery struct {
		ID                  []string `query:"id" qs:"id"`
		Name                string   `query:"name" qs:"name"`
		CrewID              []string `query:"crew_id" qs:"crew_id"`
		EventASPID          string   `query:"event_asp_id" qs:"event_asp_id"`
		Type                []string `query:"type" qs:"type"`
		EventState          []string `query:"event_state" qs:"event_state"`
		InternalASPID       string   `query:"internal_asp_id" qs:"internal_asp_id"`
		StartAt             string   `query:"start_at" qs:"start_at"`
		EndAt               string   `query:"end_at" qs:"end_at"`
		UpdatedTo           string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom         string   `query:"updated_from" qs:"updated_from"`
		MissingApplications bool     `query:"missing_applications" qs:"missing_applications"`
		OnlyApply           bool     `query:"only_apply" qs:"only_apply"`
		CreatedTo           string   `query:"created_to" qs:"created_to"`
		OrganisationId      []string `query:"organisation_id" qs:"organisation_id"`
		CreatedFrom         string   `query:"created_from" qs:"created_from"`
		FullCount           string   `query:"full_count"`
		vmdb.Query
	}
	UserExternal struct {
		FullName    string `json:"full_name" bson:"full_name"`
		DisplayName string `json:"display_name" bson:"display_name"`
		Email       string `json:"email" bson:"email"`
		Phone       string `json:"phone" bson:"phone"`
	}
	Location struct {
		Name        string   `json:"name" bson:"name"`
		Street      string   `json:"street" bson:"street"`
		Zip         string   `json:"zip" bson:"zip"`
		City        string   `json:"city" bson:"city"`
		Country     string   `json:"country" bson:"country"`
		CountryCode string   `json:"country_code" bson:"country_code"`
		Number      string   `json:"number" bson:"number"`
		Position    Position `json:"position" bson:"position"`
		PlaceID     string   `json:"place_id" bson:"place_id"`
		Sublocality string   `json:"sublocality" bson:"sublocality"`
	}
	Position struct {
		Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lng" bson:"lng"`
	}
	EventImport struct {
		Name                  string                `json:"name" bson:"name"`
		TypeOfEvent           string                `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string                `json:"additional_information" bson:"additional_information"`
		Location              Location              `json:"location" bson:"location"`
		ArtistIDs             []string              `json:"artist_ids" bson:"artist_ids"`
		Website               string                `json:"website" bson:"website"`
		OrganizerID           string                `json:"organizer_id" bson:"organizer_id"`
		StartAt               int64                 `json:"start_at" bson:"start_at"`
		EndAt                 int64                 `json:"end_at" bson:"end_at"`
		CrewID                string                `json:"crew_id" bson:"crew_id"`
		ExternalASP           UserExternal          `json:"external_asp" bson:"external_asp"`
		Application           EventApplication      `json:"application" bson:"application"`
		EventTools            EventTools            `json:"event_tools" bson:"event_tools"`
		EventState            EventState            `json:"event_state" bson:"event_state"`
		Participations        []ParticipationImport `json:"participations"`
		Modified              vmod.Modified         `json:"modified"`
	}
	EventValidate struct {
		ID                    string           `json:"id" bson:"_id"`
		Name                  string           `json:"name" bson:"name"`
		TypeOfEvent           string           `json:"type_of_event" bson:"type_of_event"`
		AdditionalInformation string           `json:"additional_information" bson:"additional_information"`
		Website               string           `json:"website" bson:"website"`
		TourID                string           `json:"tour_id" bson:"tour_id"`
		Location              Location         `json:"location" bson:"location"`
		MeetingURL            string           `json:"meeting_url" bson:"meeting_url"`
		ArtistIDs             []string         `json:"artist_ids" bson:"artist_ids"`
		Artists               []Artist         `json:"artists" bson:"artists"`
		OrganizerID           string           `json:"organizer_id" bson:"organizer_id"`
		Organizer             Organizer        `json:"organizer" bson:"organizer"`
		StartAt               int64            `json:"start_at" bson:"start_at"`
		EndAt                 int64            `json:"end_at" bson:"end_at"`
		CrewID                string           `json:"crew_id" bson:"crew_id"`
		Crew                  Crew             `json:"crew" bson:"crew"`
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		EventASP              User             `json:"event_asp" bson:"event_asp"`
		InteralASP            User             `json:"internal_asp" bson:"internal_asp"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		Participation         []Participation  `json:"participations" bson:"participations"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		Creator               User             `json:"creator" bson:"creator"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
		TakingID              string           `json:"taking_id" bson:"taking_id"`
		Taking                Taking           `json:"taking" bson:"taking"`
		EditorID              string           `json:"editor_id" bson:"editor_id"`
		Modified              vmod.Modified    `json:"modified" bson:"modified"`
	}
)

var EventCollection = "events"
var EventView = "events_view"
var PublicEventView = "events_public_view"

func (i *EventDatabase) TakingDatabase() *TakingDatabase {
	return &TakingDatabase{
		ID:           uuid.NewString(),
		Name:         i.Name,
		CrewID:       i.CrewID,
		DateOfTaking: i.EndAt,
		Type:         "automatically",
		Modified:     vmod.NewModified(),
	}
}
func (i *EventCreate) EventDatabase(token *AccessToken) *EventDatabase {
	return &EventDatabase{
		ID:                    uuid.NewString(),
		Name:                  i.Name,
		TypeOfEvent:           i.TypeOfEvent,
		AdditionalInformation: i.AdditionalInformation,
		Website:               i.Website,
		Location:              i.Location,
		MeetingURL:            i.MeetingURL,
		ArtistIDs:             i.ArtistIDs,
		OrganizerID:           i.OrganizerID,
		StartAt:               i.StartAt,
		EndAt:                 i.EndAt,
		CrewID:                i.CrewID,
		EventASPID:            i.EventASPID,
		InternalASPID:         i.InternalASPID,
		ExternalASP:           i.ExternalASP,
		Application:           i.Application,
		Applications:          EventApplications{Confirmed: 0, Rejected: 0, Requested: 0, Withdrawn: 0, Total: 0},
		EventTools:            i.EventTools,
		EventState:            i.EventState,
		OrganisationID:        i.OrganisationID,
		CreatorID:             token.ID,
		Modified:              vmod.NewModified(),
	}
}

func (i *EventDetails) FromParticipationEvent(event Event) *EventDetails {
	return &EventDetails{
		MeetingURL: event.MeetingURL,
		EventASP: EventASPPublic{
			FullName:    event.EventASP.FullName,
			DisplayName: event.EventASP.DisplayName,
			Phone:       event.EventASP.Profile.Phone,
			Email:       event.EventASP.Email,
			Mattermost:  event.EventASP.Profile.Mattermost,
		},
	}
}

func (i *EventImport) EventDatabase() *EventDatabase {
	return &EventDatabase{
		ID:                    uuid.NewString(),
		Name:                  i.Name,
		TypeOfEvent:           i.TypeOfEvent,
		AdditionalInformation: i.AdditionalInformation,
		Location:              i.Location,
		ArtistIDs:             i.ArtistIDs,
		Website:               i.Website,
		OrganizerID:           i.OrganizerID,
		StartAt:               i.StartAt,
		EndAt:                 i.EndAt,
		CrewID:                i.CrewID,
		ExternalASP:           i.ExternalASP,
		Application:           i.Application,
		EventTools:            i.EventTools,
		EventState:            i.EventState,
		Modified:              vmod.NewModified(),
	}
}

func EventPipeline(token *AccessToken) (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(UserCollection, "event_asp_id", "_id", "event_asp")
	pipe.LookupUnwind(ProfileCollection, "event_asp_id", "user_id", "event_asp.profile")
	pipe.LookupUnwind(UserCollection, "internal_asp_id", "_id", "internal_asp")
	pipe.LookupUnwind(ProfileCollection, "internal_asp_id", "user_id", "internal_asp.profile")
	pipe.LookupUnwind(UserCollection, "creator_id", "_id", "creator")
	pipe.LookupUnwind(ProfileCollection, "creator_id", "user_id", "creator.profile")
	pipe.LookupUnwind(OrganizerCollection, "organizer_id", "_id", "organizer")
	pipe.LookupUnwind(OrganisationCollection, "organisation_id", "_id", "organisation")
	if token.Roles.Validate("admin;employee;pool_employee") || token.ID == "internal" {
		pipe.Lookup(ParticipationCollection, "_id", "event_id", "participations")
	} else if token.PoolRoles.Validate(ASPEventRole) {
		pipe.LookupMatch(ParticipationEventView, "_id", "event_id", "participations", bson.D{{Key: "event.crew_id", Value: token.CrewID}})
	} else {
		pipe.LookupMatch(ParticipationEventView, "_id", "event_id", "participations", bson.D{{Key: "event.event_asp_id", Value: token.ID}})
	}
	pipe.LookupList(ArtistCollection, "artist_ids", "_id", "artists")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	return
}

func EventImportPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind(UserCollection, "event_asp_id", "_id", "event_asp")
	pipe.LookupUnwind(ProfileCollection, "event_asp_id", "user_id", "event_asp.profile")
	pipe.LookupUnwind(UserCollection, "internal_asp_id", "_id", "internal_asp")
	pipe.LookupUnwind(ProfileCollection, "internal_asp_id", "user_id", "internal_asp.profile")
	pipe.LookupUnwind(UserCollection, "creator_id", "_id", "creator")
	pipe.LookupUnwind(ProfileCollection, "creator_id", "user_id", "creator.profile")
	pipe.Lookup(ParticipationCollection, "_id", "event_id", "participations")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	return
}

func EventPipelinePublic() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.Lookup(ParticipationCollection, "_id", "event_id", "participations")
	pipe.LookupUnwind(OrganizerCollection, "organizer_id", "_id", "organizer")
	pipe.LookupList(ArtistCollection, "artist_ids", "_id", "artists")
	pipe.LookupUnwind(CrewCollection, "crew_id", "_id", "crew")
	pipe.LookupUnwind(OrganisationCollection, "organisation_id", "_id", "organisation")
	return
}

func EventRolePipeline() *vmdb.Pipeline {
	pipe := vmdb.NewPipeline()
	pipe.LookupUnwindMatch(PoolRoleCollection, "user_id", "user_id", "user", bson.D{{Key: "name", Value: "operation"}})
	return pipe
}

func EventPermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		return vcago.NewPermissionDenied(EventCollection)
	}
	return
}

func EventDeletePermission(token *AccessToken) (err error) {
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		return vcago.NewPermissionDenied(EventCollection)
	}
	return
}

func (i *EventParam) EventSyncPermission(token *AccessToken) (err error) {
	if !token.Roles.Validate("admin") {
		return vcago.NewPermissionDenied(EventCollection)
	}
	return
}

func (i *EventParam) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *EventDatabase) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *EventParam) PermittedDeleteFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew_id", token.CrewID)
		filter.EqualStringList("event_state.state", []string{"created", "requested"})
	}
	return filter.Bson()
}
func (i *EventUpdate) Match() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *EventUpdate) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualString("event_asp_id", token.ID)
	} else if !token.Roles.Validate("admin;employee;pool_employee") {
		crewMatch := vmdb.NewFilter()
		crewMatch.EqualString("crew_id", token.CrewID)
		eventAspMatch := vmdb.NewFilter()
		eventAspMatch.EqualString("event_asp_id", token.ID)
		filter.Append(bson.E{Key: "$or", Value: bson.A{eventAspMatch.Bson(), crewMatch.Bson()}})
	}
	return filter.Bson()
}

func (i *Event) FilterCrew() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualBool("confirmed", "true")
	filter.EqualString("crew.crew_id", i.CrewID)
	filter.ElemMatchList("pool_roles", "name", []string{"education", "network", "operation"})
	return filter.Bson()
}

func (i *EventParam) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualString("event_asp_id", token.ID)
	} else if !token.Roles.Validate("admin;employee;pool_employee") {
		crewMatch := vmdb.NewFilter()
		crewMatch.EqualString("crew_id", token.CrewID)
		eventAspMatch := vmdb.NewFilter()
		eventAspMatch.EqualString("event_asp_id", token.ID)
		filter.Append(bson.E{Key: "$or", Value: bson.A{eventAspMatch.Bson(), crewMatch.Bson()}})
	}

	return filter.Bson()
}

func (i *EventParam) FilterID() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	return filter.Bson()
}

func (i *EventQuery) PublicFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualStringList("_id", i.ID)
	filter.LikeString("name", i.Name)
	filter.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	filter.EqualStringList("type_of_event", i.Type)
	filter.EqualStringList("organisation_id", i.OrganisationId)
	filter.EqualStringList("crew_id", i.CrewID)
	filter.GteInt64("start_at", i.StartAt)
	filter.LteInt64("end_at", i.EndAt)
	if i.OnlyApply {
		filter.LteInt64("application.start_date", strconv.FormatInt(time.Now().Unix(), 10))
		filter.GteInt64("application.end_date", strconv.FormatInt(time.Now().Unix(), 10))
	}
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	filter.SearchString([]string{"_id", "name", "crew.name", "artists.name", "location.name"}, i.Search)

	return filter.Bson()
}

func (i EventQuery) Sort() bson.D {
	sort := vmdb.NewSort()
	sort.Add(i.SortField, i.SortDirection)
	return sort.Bson()
}

func (i *EventParam) PublicFilter() bson.D {
	filter := vmdb.NewFilter()
	filter.EqualString("_id", i.ID)
	filter.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	return filter.Bson()
}

func (i *EventQuery) FilterAsp(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if token.PoolRoles.Validate(ASPEventRole) {
		status := bson.A{}
		status = append(status, bson.D{{Key: "event_asp_id", Value: token.ID}})
		status = append(status, bson.D{{Key: "crew_id", Value: token.CrewID}})
		filter.Append(bson.E{Key: "$or", Value: status})
	} else {
		filter.EqualString("event_asp_id", token.ID)
	}
	filter.GteInt64("end_at", fmt.Sprint(time.Now().AddDate(0, -6, 0).Unix()))
	filter.EqualStringList("event_state.state", []string{"draft", "published", "finished"})
	return filter.Bson()
}

func (i *EventQuery) PermittedFilter(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	} else if !token.Roles.Validate("admin;employee;pool_employee") {
		noCrewMatch := vmdb.NewFilter()
		crewMatch := vmdb.NewFilter()
		crewMatch.EqualString("crew_id", token.CrewID)
		noCrewMatch.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
		filter.Append(bson.E{Key: "$or", Value: bson.A{noCrewMatch.Bson(), crewMatch.Bson()}})
	}
	if i.OnlyApply {
		filter.LteInt64("application.start_date", strconv.FormatInt(time.Now().Unix(), 10))
		filter.GteInt64("application.end_date", strconv.FormatInt(time.Now().Unix(), 10))
	}
	if i.MissingApplications {
		filter.Append(bson.E{Key: "$expr", Value: bson.D{{Key: "$gt", Value: bson.A{"$application.supporter_count", "$applications.confirmed"}}}})
	}
	filter.GteInt64("start_at", i.StartAt)
	filter.LteInt64("end_at", i.EndAt)
	filter.EqualStringList("_id", i.ID)
	filter.LikeString("name", i.Name)
	filter.EqualStringList("type_of_event", i.Type)
	filter.EqualString("internal_asp_id", i.InternalASPID)
	filter.EqualString("event_asp_id", i.EventASPID)
	filter.EqualStringList("event_state.state", i.EventState)
	filter.EqualStringList("crew_id", i.CrewID)
	filter.EqualStringList("organisation_id", i.OrganisationId)
	filter.GteInt64("modified.updated", i.UpdatedFrom)
	filter.GteInt64("modified.created", i.CreatedFrom)
	filter.LteInt64("modified.updated", i.UpdatedTo)
	filter.LteInt64("modified.created", i.CreatedTo)
	filter.SearchString([]string{"_id", "name", "crew.name", "artists.name", "location.name"}, i.Search)
	return filter.Bson()
}

func (i *EventQuery) FilterEmailEvents(token *AccessToken) bson.D {
	filter := vmdb.NewFilter()
	if !(token.Roles.Validate("admin;employee;pool_employee") || token.PoolRoles.Validate(ASPEventRole)) {
		filter.EqualString("event_asp_id", token.ID)
		filter.EqualString("crew_id", token.CrewID)
	} else if !token.Roles.Validate("admin;employee;pool_employee") {
		filter.EqualString("crew_id", token.CrewID)
	}

	return filter.Bson()
}

func (i *Event) ToContent() *vmod.Content {
	content := &vmod.Content{
		Fields: make(map[string]interface{}),
	}
	content.Fields["Event"] = i
	return content
}

func (i *EventUpdate) EventStateValidation(token *AccessToken, event *EventValidate) (err error) {
	if i.EventState.State == "canceled" && (event.EventState.State == "finished" || event.EventState.State == "closed") {
		return vcago.NewBadRequest(EventCollection, "event can not be canceled, it is already "+event.EventState.State, i)
	}
	if !token.Roles.Validate("admin;employee;pool_employee") && (event.EventState.State == "finished" || event.EventState.State == "closed") {
		return vcago.NewBadRequest(EventCollection, "event can not be updated, it is already "+event.EventState.State, i)
	}
	if event.Taking.Money.Amount != 0 {
		return vcago.NewBadRequest(EventCollection, "taking_failure", nil)
	}
	return
}

func (i *Event) GetLocation() string {
	if i.Location.PlaceID != "" {
		return fmt.Sprintf("%v, %v %v", i.Location.Name, i.Location.Zip, i.Location.City)
	}
	return fmt.Sprintf("Online: %v", i.MeetingURL)
}

func (i *EventPublic) GetLocation() string {
	if i.Location.PlaceID != "" {
		return fmt.Sprintf("%v, %v %v", i.Location.Name, i.Location.Zip, i.Location.City)
	}
	return "Online"
}
