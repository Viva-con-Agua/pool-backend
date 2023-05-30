package models

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/Viva-con-Agua/vcapool"
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
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
	}
	EventDatabase struct {
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
		TakingID              string           `json:"taking_id" bson:"taking_id"`
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		CreatorID             string           `json:"creator_id" bson:"creator_id"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified    `json:"modified" bson:"modified"`
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
	//EventState represents the state of an event.
	EventState struct {
		State                string `json:"state" bson:"state"`
		CrewConfirmation     string `json:"crew_confirmation" bson:"crew_confirmation"`
		InternalConfirmation string `json:"internal_confirmation" bson:"internal_confirmation"`
		TakingID             string `json:"taking_id" bson:"taking_id"`
		DepositID            string `json:"deposit_id" bson:"deposit_id"`
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
		StartAt               int64                  `json:"start_at" bson:"start_at"`
		EndAt                 int64                  `json:"end_at" bson:"end_at"`
		CrewID                string                 `json:"crew_id" bson:"crew_id"`
		Crew                  Crew                   `json:"crew" bson:"crew"`
		EventASPID            string                 `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string                 `json:"internal_asp_id" bson:"internal_asp_id"`
		Application           EventApplication       `json:"application" bson:"application"`
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
		Participation         []ParticipationMinimal `json:"participations" bson:"participations"`
		EventTools            EventTools             `json:"event_tools" bson:"event_tools"`
		EventState            EventState             `json:"event_state" bson:"event_state"`
		Modified              vmod.Modified          `json:"modified" bson:"modified"`
	}
	Event struct {
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
		EditorID              string           `json:"editor_id" bson:"editor_id"`
		Modified              vmod.Modified    `json:"modified" bson:"modified"`
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
		Crew                  Crew             `json:"crew" bson:"crew"`
		EventASPID            string           `json:"event_asp_id" bson:"event_asp_id"`
		InternalASPID         string           `json:"internal_asp_id" bson:"internal_asp_id"`
		ExternalASP           UserExternal     `json:"external_asp" bson:"external_asp"`
		Application           EventApplication `json:"application" bson:"application"`
		EventTools            EventTools       `json:"event_tools" bson:"event_tools"`
		EventState            EventState       `json:"event_state" bson:"event_state"`
	}
	EventParam struct {
		ID string `param:"id"`
	}

	EventQuery struct {
		ID            []string `query:"id" qs:"id"`
		Name          string   `query:"name" qs:"name"`
		CrewID        string   `query:"crew_id" qs:"crew_id"`
		EventASPID    string   `query:"event_asp_id" qs:"event_asp_id"`
		EventState    []string `query:"event_state" qs:"event_state"`
		InternalASPID string   `query:"internal_asp_id" qs:"internal_asp_id"`
		UpdatedTo     string   `query:"updated_to" qs:"updated_to"`
		UpdatedFrom   string   `query:"updated_from" qs:"updated_from"`
		CreatedTo     string   `query:"created_to" qs:"created_to"`
		CreatedFrom   string   `query:"created_from" qs:"created_from"`
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
)

func (i *EventCreate) EventDatabase(token *vcapool.AccessToken) *EventDatabase {
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
		EventTools:            i.EventTools,
		EventState:            i.EventState,
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

func EventPipeline(token *vcapool.AccessToken) (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind("users", "event_asp_id", "_id", "event_asp")
	pipe.LookupUnwind("profiles", "event_asp_id", "user_id", "event_asp.profile")
	pipe.LookupUnwind("users", "internal_asp_id", "_id", "internal_asp")
	pipe.LookupUnwind("profiles", "internal_asp_id", "user_id", "internal_asp.profile")
	pipe.LookupUnwind("users", "creator_id", "_id", "creator")
	pipe.LookupUnwind("profiles", "creator_id", "user_id", "creator.profile")
	pipe.LookupUnwind("organizers", "organizer_id", "_id", "organizer")
	if token.Roles.Validate("employee;admin") {
		pipe.Lookup("participations", "_id", "event_id", "participations")
	} else if token.PoolRoles.Validate("network;operation;education") {
		pipe.LookupMatch("participations_event", "_id", "event_id", "participations", bson.D{{Key: "event.crew_id", Value: token.CrewID}})
	} else {
		pipe.LookupMatch("participations_event", "_id", "event_id", "participations", bson.D{{Key: "event.event_asp_id", Value: token.ID}})
	}
	pipe.LookupList("artists", "artist_ids", "_id", "artists")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	return
}

func EventImportPipeline() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.LookupUnwind("users", "event_asp_id", "_id", "event_asp")
	pipe.LookupUnwind("profiles", "event_asp_id", "user_id", "event_asp.profile")
	pipe.LookupUnwind("users", "internal_asp_id", "_id", "internal_asp")
	pipe.LookupUnwind("profiles", "internal_asp_id", "user_id", "internal_asp.profile")
	pipe.LookupUnwind("users", "creator_id", "_id", "creator")
	pipe.LookupUnwind("profiles", "creator_id", "user_id", "creator.profile")
	pipe.Lookup("participations", "_id", "event_id", "participations")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	return
}

func EventPipelinePublic() (pipe *vmdb.Pipeline) {
	pipe = vmdb.NewPipeline()
	pipe.Lookup("participations", "_id", "event_id", "participations")
	pipe.LookupUnwind("organizers", "organizer_id", "_id", "organizer")
	pipe.LookupList("artists", "artist_ids", "_id", "artists")
	pipe.LookupUnwind("crews", "crew_id", "_id", "crew")
	return
}

func EventPermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied("crew", nil)
	}
	return
}

func EventDeletePermission(token *vcapool.AccessToken) (err error) {
	if !token.Roles.Validate("employee;admin") {
		return vcago.NewPermissionDenied("event", nil)
	}
	return
}

func (i *EventDatabase) Match() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	return match.Bson()
}

func (i *EventParam) Match() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	return match.Bson()
}
func (i *EventUpdate) Match() bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	return match.Bson()
}

func (i *EventUpdate) PermittedFilter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)

	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualString("event_asp_id", token.ID)
		match.EqualString("crew_id", token.CrewID)
	} else if !token.Roles.Validate("employee;admin") {
		match.EqualString("crew_id", token.CrewID)
	}

	return match.Bson()
}

func (i *EventParam) PermittedFilter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	match.EqualString("_id", i.ID)
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualString("event_asp_id", token.ID)
		match.EqualString("crew_id", token.CrewID)
		match.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	} else if !token.Roles.Validate("employee;admin") {
		match.EqualString("crew_id", token.CrewID)
	}
	return match.Bson()
}

func (i *EventQuery) Match() bson.D {
	match := vmdb.NewFilter()
	match.EqualStringList("_id", i.ID)
	match.LikeString("name", i.Name)
	match.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	match.EqualString("crew_id", i.CrewID)
	match.GteInt64("modified.updated", i.UpdatedFrom)
	match.GteInt64("modified.created", i.CreatedFrom)
	match.LteInt64("modified.updated", i.UpdatedTo)
	match.LteInt64("modified.created", i.CreatedTo)
	return match.Bson()
}

func (i *EventQuery) FilterAsp(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	match.EqualString("event_asp_id", token.ID)
	return match.Bson()
}

func (i *EventQuery) PermittedFilter(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
	} else if !token.Roles.Validate("employee;admin") {
		noCrewMatch := vmdb.NewFilter()
		crewMatch := vmdb.NewFilter()
		crewMatch.EqualString("crew_id", token.CrewID)
		noCrewMatch.EqualStringList("event_state.state", []string{"published", "finished", "closed"})
		match.Append(bson.E{Key: "$or", Value: bson.A{noCrewMatch.Bson(), crewMatch.Bson()}})
	}
	match.EqualStringList("_id", i.ID)
	match.LikeString("name", i.Name)
	match.EqualString("internal_asp_id", i.InternalASPID)
	match.EqualString("event_asp_id", i.EventASPID)
	match.EqualStringList("event_state.state", i.EventState)
	match.EqualString("crew_id", i.CrewID)
	match.GteInt64("modified.updated", i.UpdatedFrom)
	match.GteInt64("modified.created", i.CreatedFrom)
	match.LteInt64("modified.updated", i.UpdatedTo)
	match.LteInt64("modified.created", i.CreatedTo)
	return match.Bson()
}

func (i *EventQuery) FilterEmailEvents(token *vcapool.AccessToken) bson.D {
	match := vmdb.NewFilter()
	if !(token.Roles.Validate("employee;admin") || token.PoolRoles.Validate("network;operation;education")) {
		match.EqualString("event_asp_id", token.ID)
		match.EqualString("crew_id", token.CrewID)
	} else if !token.Roles.Validate("employee;admin") {
		match.EqualString("crew_id", token.CrewID)
	}

	return match.Bson()
}
