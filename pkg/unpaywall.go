package pkg

// See the data description at https://unpaywall.org/data-format

// DOI represents the object for a DOI-assigned resource, including metadata and information about its OA status.
type DOI struct {
	// BestOA is the best OA Location Object for the DOI.
	BestOA *OA `json:"best_oa_location,omitempty"`
	// DataStandard is the data collection approaches used for the resource.
	DataStandard int `json:"data_standard"`
	// DOI is the identifier for the resource.
	DOI string `json:"doi"`
	// DOIURL is the DOI in hyperlink form.
	DOIURL string `json:"doi_url"`
	// Genre is the type of resource.
	Genre string `json:"genre,omitempty"`
	// IsParatext represents whether the item is an ancillary part of a journal.
	IsParatext bool `json:"is_paratext"`
	// IsOA represents whether there is an OA copy of the resource.
	IsOA bool `json:"is_oa"`
	// JournalIsInDOAJ represents whether the resource is published in a DOAJ-indexed journal.
	JournalIsInDOAJ bool `json:"journal_is_in_doaj"`
	// JournalIsOA represents whether the resource is published in a completely OA journal.
	JournalIsOA bool `json:"journal_is_oa"`
	// JournalISSNs represents any ISSNs assigned to the journal publishing the resource.
	JournalISSNs string `json:"journal_issns,omitempty"`
	// JournalISSNL represents a single ISSN for the journal publishing the resource.
	JournalISSNL string `json:"journal_issn_l,omitempty"`
	// JournalName represents the name of the journal publishing the resource.
	JournalName string `json:"journal_name,omitempty"`
	// OALocations represents the list of all OA Location objects associated with the resource.
	OALocations []OA `json:"oa_locations,omitempty"`
	// OALocationsEmbargoed represents the list of OA Location objects that are not yet available.
	OALocationsEmbargoed []OA `json:"oa_locations_embargoed,omitempty"`
	// FirstOALocation represents the OA Location Object with the earliest oa_date.
	FirstOALocation *OA `json:"first_oa_location,omitempty"`
	// OAStatus represents the OA status of the resource.
	OAStatus string `json:"oa_status"`
	// HasRepositoryCopy represents whether there is a copy of the resource in a repository.
	HasRepositoryCopy bool `json:"has_repository_copy"`
	// PublishedDate represents the date the resource was published.
	PublishedDate string `json:"published_date,omitempty"`
}

// OA describes a location of an Open Access article
type OA struct {
	// Evidence describes how the OA location was found
	Evidence string `json:"evidence"`
	// HostType describes the type of host serving the OA location
	HostType string `json:"host_type"`
	// IsBest indicates if this is the best OA location for the article
	IsBest bool `json:"is_best"`
	// License describes the license under which the OA article is published
	License *string `json:"license,omitempty"`
	// OADate describes when the article first became available at this location
	OADate *string `json:"oa_date,omitempty"`
	// PMHID is the OAI-PMH endpoint where the OA location was found
	PMHID *string `json:"pmh_id,omitempty"`
	// Updated is the time when the data for this location was last updated
	Updated string `json:"updated"`
	// URL is the URL for the PDF version of the OA article or the landing page URL if no PDF is available
	URL string `json:"url"`
	// URLForLandingPage is the URL for the landing page describing the OA article
	URLForLandingPage string `json:"url_for_landing_page"`
	// URLForPDF is the URL for the PDF version of the OA article
	URLForPDF *string `json:"url_for_pdf,omitempty"`
	// Version is the content version accessible at this location
	Version string `json:"version"`
}
