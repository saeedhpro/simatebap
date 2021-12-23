package medicalHistory

import (
	"fmt"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type MedicalHistoryStruct struct {
	UserID                  *int64 `json:"user_id"`
	AdenoidTonsileReduction string `json:"adenoid_tonsile_reduction"`
	MedicalCondition        string `json:"medical_condition"`
	ConsumableMedicine      string `json:"consumable_medicine"`
	GeneralHealth           string `json:"general_health"`
	UnderPhysicianCare      string `json:"under_physician_care"`
	AccidentToHead          string `json:"accident_to_Head"`
	Operations              string `json:"operations"`
	ChiefComplaint          string `json:"chief_complaint"`
	PreviousOrthodontic     string `json:"previous_orthodontic"`
	OralHygiene             string `json:"oral_hygiene"`
	Frontal                 string `json:"frontal"`
	Profile                 string `json:"profile"`
	TeethPresent            string `json:"teeth_present"`
	UnErupted               string `json:"un_erupted"`
	IeMissing               string `json:"ie_missing"`
	IeExtracted             string `json:"ie_extracted"`
	IeImpacted              string `json:"ie_impacted"`
	IeSupernumerary         string `json:"ie_supernumerary"`
	IeCaries                string `json:"ie_caries"`
	IeRct                   string `json:"ie_rct"`
	IeAnomalies             string `json:"ie_anomalies"`
	LeftMolar               string `json:"left_molar"`
	RightMolar              string `json:"right_molar"`
	LeftCanine              string `json:"left_canine"`
	RightCanine             string `json:"right_canine"`
	Overjet                 string `json:"overjet"`
	Overbite                string `json:"overbite"`
	Crossbite               string `json:"crossbite"`
	CrowdingMd              string `json:"crowding_md"`
	CrowdingMx              string `json:"crowding_mx"`
	SpacingMx               string `json:"spacing_mx"`
	SpacingMd               string `json:"spacing_md"`
	Diagnosis               string `json:"diagnosis"`
	TreatmentPlan           string `json:"treatment_plan"`
	LengthActiveTreatment   string `json:"length_active_treatment"`
	Retention               string `json:"retention"`
	Prognosis               string `json:"prognosis"`
	Consultation            string `json:"consultation"`
	Charging                string `json:"charging"`
}

type CreateMedicalHistoryStruct struct {
	UserID                  *int64 `json:"user_id"`
	AdenoidTonsileReduction string `json:"adenoid_tonsile_reduction"`
	MedicalCondition        string `json:"medical_condition"`
	ConsumableMedicine      string `json:"consumable_medicine"`
	GeneralHealth           string `json:"general_health"`
	UnderPhysicianCare      string `json:"under_physician_care"`
	AccidentToHead          string `json:"accident_to_Head"`
	Operations              string `json:"operations"`
	ChiefComplaint          string `json:"chief_complaint"`
	PreviousOrthodontic     string `json:"previous_orthodontic"`
	OralHygiene             string `json:"oral_hygiene"`
	Frontal                 string `json:"frontal"`
	Profile                 string `json:"profile"`
	TeethPresent            string `json:"teeth_present"`
	UnErupted               string `json:"un_erupted"`
	IeMissing               string `json:"ie_missing"`
	IeExtracted             string `json:"ie_extracted"`
	IeImpacted              string `json:"ie_impacted"`
	IeSupernumerary         string `json:"ie_supernumerary"`
	IeCaries                string `json:"ie_caries"`
	IeRct                   string `json:"ie_rct"`
	IeAnomalies             string `json:"ie_anomalies"`
	LeftMolar               string `json:"left_molar"`
	RightMolar              string `json:"right_molar"`
	LeftCanine              string `json:"left_canine"`
	RightCanine             string `json:"right_canine"`
	Overjet                 string `json:"overjet"`
	Overbite                string `json:"overbite"`
	Crossbite               string `json:"crossbite"`
	CrowdingMd              string `json:"crowding_md"`
	CrowdingMx              string `json:"crowding_mx"`
	SpacingMx               string `json:"spacing_mx"`
	SpacingMd               string `json:"spacing_md"`
	Diagnosis               string `json:"diagnosis"`
	TreatmentPlan           string `json:"treatment_plan"`
	LengthActiveTreatment   string `json:"length_active_treatment"`
	Retention               string `json:"retention"`
	Prognosis               string `json:"prognosis"`
	Consultation            string `json:"consultation"`
	Charging                string `json:"charging"`
}

func GetMedicalHistory(userID string) (*MedicalHistoryStruct, error) {
	history := MedicalHistoryStruct{}
	query := "SELECT " +
		"ifnull(user_id, 0) user_id," +
		"ifnull(adenoid_tonsile_reduction, '') adenoid_tonsile_reduction," +
		"ifnull(medical_condition, '') medical_condition," +
		"ifnull(consumable_medicine, '') consumable_medicine," +
		"ifnull(general_health, '') general_health," +
		"ifnull(under_physician_care, '') under_physician_care," +
		"ifnull(accident_to_Head, '') accident_to_Head," +
		"ifnull(operations, '') operations," +
		"ifnull(chief_complaint, '') chief_complaint," +
		"ifnull(previous_orthodontic, '') previous_orthodontic," +
		"ifnull(oral_hygiene, '') oral_hygiene," +
		"ifnull(frontal, '') frontal," +
		"ifnull(profile, '') 'profile'," +
		"ifnull(teeth_present, '') teeth_present," +
		"ifnull(un_erupted, '') un_erupted," +
		"ifnull(ie_missing, '') ie_missing," +
		"ifnull(ie_extracted, '') ie_extracted," +
		"ifnull(ie_impacted, '') ie_impacted," +
		"ifnull(ie_supernumerary, '') ie_supernumerary," +
		"ifnull(ie_caries, '') ie_caries," +
		"ifnull(ie_rct, '') ie_rct," +
		"ifnull(ie_anomalies, '') ie_anomalies," +
		"ifnull(left_molar, '') left_molar," +
		"ifnull(right_molar, '') right_molar," +
		"ifnull(left_canine, '') left_canine," +
		"ifnull(right_canine, '') right_canine," +
		"ifnull(overjet, '') overjet," +
		"ifnull(overbite, '') overbite," +
		"ifnull(crossbite, '') crossbite," +
		"ifnull(crowding_md, '') crowding_md," +
		"ifnull(crowding_mx, '') crowding_mx," +
		"ifnull(spacing_mx, '') spacing_mx," +
		"ifnull(spacing_md, '') spacing_md," +
		"ifnull(diagnosis, '') diagnosis," +
		"ifnull(treatment_plan, '') treatment_plan," +
		"ifnull(length_active_treatment, '') length_active_treatment," +
		"ifnull(retention, '') retention," +
		"ifnull(prognosis, '') prognosis," +
		"ifnull(consultation, '') consultation," +
		"ifnull(charging, '') charging" +
		" FROM medical_history_orthodontics WHERE user_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return nil, err
	}
	row := stmt.QueryRow(userID)
	err = row.Scan(
		&history.UserID,
		&history.AdenoidTonsileReduction,
		&history.MedicalCondition,
		&history.ConsumableMedicine,
		&history.GeneralHealth,
		&history.UnderPhysicianCare,
		&history.AccidentToHead,
		&history.Operations,
		&history.ChiefComplaint,
		&history.PreviousOrthodontic,
		&history.OralHygiene,
		&history.Frontal,
		&history.Profile,
		&history.TeethPresent,
		&history.UnErupted,
		&history.IeMissing,
		&history.IeExtracted,
		&history.IeImpacted,
		&history.IeSupernumerary,
		&history.IeCaries,
		&history.IeRct,
		&history.IeAnomalies,
		&history.LeftMolar,
		&history.RightMolar,
		&history.LeftCanine,
		&history.RightCanine,
		&history.Overjet,
		&history.Overbite,
		&history.Crossbite,
		&history.CrowdingMd,
		&history.CrowdingMx,
		&history.SpacingMx,
		&history.SpacingMd,
		&history.Diagnosis,
		&history.TreatmentPlan,
		&history.LengthActiveTreatment,
		&history.Retention,
		&history.Prognosis,
		&history.Consultation,
		&history.Charging,
	)
	if err != nil {
		log.Println(err.Error(), " :err: ")
		return nil, err
	}
	return &history, nil
}

func CreateMedicalHistory(request CreateMedicalHistoryStruct) error {
	h, err := GetMedicalHistory(fmt.Sprintf("%d", request.UserID))
	if err != nil {
		log.Println(err.Error())
	}
	if h == nil {
		query := "UPDATE `medical_history_orthodontics` SET" +
			"`adenoid_tonsile_reduction` = ? ,`medical_condition` = ? ," +
			"`consumable_medicine`= ? ,`general_health`= ? ,`under_physician_care` = ? ," +
			"`accident_to_Head`= ? ,`operations`= ?,`chief_complaint`= ? ," +
			"`previous_orthodontic`= ? ,`oral_hygiene`= ? ,`frontal`= ? ," +
			"`profile`= ? ,`teeth_present`= ? ,`un_erupted`= ? ,`ie_missing`= ? ," +
			"`ie_extracted`= ?,`ie_impacted`= ?,`ie_supernumerary`= ?,`ie_caries`=?," +
			"`ie_rct`=?,`ie_anomalies`=?,`left_molar`=?,`right_molar`=?," +
			"`left_canine`=?,`right_canine`=?,`overjet`=?,`overbite`=?," +
			"`crossbite`=?,`crowding_md`=?,`crowding_mx`= ?,`spacing_mx`=?,`spacing_md`=?," +
			"`diagnosis`=?,`treatment_plan`=?,`length_active_treatment`= ?,`retention`= ?," +
			"`prognosis`=?,`consultation`=?,`charging`= ? " +
			" WHERE user_id = ?"
		stmt, err := repository.DBS.MysqlDb.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			&request.AdenoidTonsileReduction,
			&request.MedicalCondition,
			&request.ConsumableMedicine,
			&request.GeneralHealth,
			&request.UnderPhysicianCare,
			&request.AccidentToHead,
			&request.Operations,
			&request.ChiefComplaint,
			&request.PreviousOrthodontic,
			&request.OralHygiene,
			&request.Frontal,
			&request.Profile,
			&request.TeethPresent,
			&request.UnErupted,
			&request.IeMissing,
			&request.IeExtracted,
			&request.IeImpacted,
			&request.IeSupernumerary,
			&request.IeCaries,
			&request.IeRct,
			&request.IeAnomalies,
			&request.LeftMolar,
			&request.RightMolar,
			&request.LeftCanine,
			&request.RightCanine,
			&request.Overjet,
			&request.Overbite,
			&request.Crossbite,
			&request.CrowdingMd,
			&request.CrowdingMx,
			&request.SpacingMx,
			&request.SpacingMd,
			&request.Diagnosis,
			&request.TreatmentPlan,
			&request.LengthActiveTreatment,
			&request.Retention,
			&request.Prognosis,
			&request.Consultation,
			&request.Charging,
			&request.UserID,
		)
		if err != nil {
			return err
		}
		return nil
	} else {
		query := "INSERT INTO `medical_history_orthodontics`(`user_id`, `adenoid_tonsile_reduction`, `medical_condition`, `consumable_medicine`, `general_health`, `under_physician_care`, `accident_to_Head`, `operations`, `chief_complaint`, `previous_orthodontic`, `oral_hygiene`, `frontal`, `profile`, `teeth_present`, `un_erupted`, `ie_missing`, `ie_extracted`, `ie_impacted`, `ie_supernumerary`, `ie_caries`, `ie_rct`, `ie_anomalies`, `left_molar`, `right_molar`, `left_canine`, `right_canine`, `overjet`, `overbite`, `crossbite`, `crowding_md`, `crowding_mx`, `spacing_mx`, `spacing_md`, `diagnosis`, `treatment_plan`, `length_active_treatment`, `retention`, `prognosis`, `consultation`, `charging`)" +
			" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) "
		stmt, err := repository.DBS.MysqlDb.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(
			&request.UserID,
			&request.AdenoidTonsileReduction,
			&request.MedicalCondition,
			&request.ConsumableMedicine,
			&request.GeneralHealth,
			&request.UnderPhysicianCare,
			&request.AccidentToHead,
			&request.Operations,
			&request.ChiefComplaint,
			&request.PreviousOrthodontic,
			&request.OralHygiene,
			&request.Frontal,
			&request.Profile,
			&request.TeethPresent,
			&request.UnErupted,
			&request.IeMissing,
			&request.IeExtracted,
			&request.IeImpacted,
			&request.IeSupernumerary,
			&request.IeCaries,
			&request.IeRct,
			&request.IeAnomalies,
			&request.LeftMolar,
			&request.RightMolar,
			&request.LeftCanine,
			&request.RightCanine,
			&request.Overjet,
			&request.Overbite,
			&request.Crossbite,
			&request.CrowdingMd,
			&request.CrowdingMx,
			&request.SpacingMx,
			&request.SpacingMd,
			&request.Diagnosis,
			&request.TreatmentPlan,
			&request.LengthActiveTreatment,
			&request.Retention,
			&request.Prognosis,
			&request.Consultation,
			&request.Charging,
		)
		if err != nil {
			return err
		}
		return nil
	}
}
