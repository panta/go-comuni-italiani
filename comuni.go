package comuni

type Comune struct {
	CodiceRegione                  string `json:"codice_regione" csv:"Codice Regione"`
	CodiceUnitaTerritoriale        string `json:"codice_unita_territoriale" csv:"Codice dell'Unità territoriale sovracomunale..."`
	CodiceProvincia                string `json:"codice_provincia" csv:"Codice Provincia (Storico)..."`
	ProgressivoComune              string `json:"progressivo_comune" csv:"Progressivo del Comune..."`
	CodiceComune                   string `json:"codice_comune" csv:"Codice Comune formato alfanumerico"`
	Denominazione                  string `json:"denominazione" csv:"Denominazione (Italiana e straniera)"`
	DenominazioneIta               string `json:"denominazione_ita" csv:"Denominazione in italiano"`
	DenominazioneNonIta            string `json:"denominazione_non_ita" csv:"Denominazione altra lingua"`
	CodiceRipartizioneGeo          string `json:"codice_ripartizione_geo" csv:"Codice Ripartizione Geografica"`
	RipartizioneGeo                string `json:"ripartizione_geo" csv:"Ripartizione geografica"`
	DenominazioneRegione           string `json:"denominazione_regione" csv:"Denominazione Regione"`
	DenominazioneUnitaTerritoriale string `json:"denominazione_unita_territoriale" csv:"Denominazione dell'Unità territoriale sovracomunale..."`
	TipologiaUnitaTerritoriale     string `json:"tipologia_unita_territoriale" csv:"Tipologia di Unità territoriale sovracomunale"`
	ComuneCapoluogo                bool   `json:"comune_capoluogo" csv:"Flag Comune capoluogo di provincia/città..."`
	SiglaAutomobilistica           string `json:"sigla_automobilistica" csv:"Sigla automobilistica"`
	CodiceComuneNumerico           int32  `json:"codice_comune_numerico" csv:"Codice Comune formato numerico"`
	CodiceComuneNumerico110        int32  `json:"codice_comune_numerico_110" csv:"Codice Comune numerico con 110 province..."`
	CodiceComuneNumerico107        int32  `json:"codice_comune_numerico_107" csv:"Codice Comune numerico con 107 province..."`
	CodiceComuneNumerico103        int32  `json:"codice_comune_numerico_103" csv:"Codice Comune numerico con 103 province..."`
	Codice_NUTS1_2010              string `json:"codice_nuts1_2010" csv:"Codice NUTS1 2010"`
	Codice_NUTS2_2010              string `json:"codice_nuts2_2010" csv:"Codice NUTS2 2010..."`
	Codice_NUTS3_2010              string `json:"codice_nuts3_2010" csv:"Codice NUTS3 2010"`
	Codice_NUTS1_2021              string `json:"codice_nuts1_2021" csv:"Codice NUTS1 2021"`
	Codice_NUTS2_2021              string `json:"codice_nuts2_2021" csv:"Codice NUTS2 2021..."`
	Codice_NUTS3_2021              string `json:"codice_nuts3_2021" csv:"Codice NUTS3 2021"`
}
