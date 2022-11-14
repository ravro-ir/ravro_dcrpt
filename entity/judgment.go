package entity

type RavroVulnerability struct {
	Name    string `json:"name"`
	Define  string `json:"define"`
	Fix     string `json:"fix"`
	Writeup string `json:"writeup"`
}

type JudgeCvss struct {
	Value  string `json:"value"`
	Rating string `json:"rating"`
}

type Judgment struct {
	Score         int    `json:"score"`
	Reward        int    `json:"reward"`
	Description   string `json:"description"`
	Cvss          JudgeCvss
	Vulnerability RavroVulnerability
}
