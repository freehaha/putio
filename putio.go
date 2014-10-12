package putio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// fix problems with json not handling null as a value

const (
	BaseUrl = "https://api.put.io/v2/"
)

var oauthparam = "?oauth_token="
var oathtoken string

type MP4 struct {
	Status       string
	Stream_url   string
	Download_url string
	Size         int64
	Percent_done int
}

type File struct {
	Is_shared          bool   `json:"is_shared"`
	Name               string `json:"name"`
	Screenshot         string `json:"screenshot"` // returns url to image
	Created_at         string `json:"created_at"` // in iso8601 format
	Opensubtitles_hash string `json:"opensubtitles_hash"`
	Parent_id          int64  `json:"parent_id"` // parent folder id
	Is_mp4_available   bool   `json:"is_mp4_available"`
	Content_type       string `json:"content_type"`
	Crc32              string `json:"crc32"`
	Icon               string `json:"icon"` // returns url to screenshot image in icon size
	Id                 int64  `json:"id"`
	Size               int64  `json:"size"`
}

type Files struct {
	Files  []File // for multi file results
	File   File   // for single file result like files/id
	Mp4    MP4    // for mp4 streaming results
	Status string
	Parent File
	Next   string
}

type Transfer struct {
	Uploaded        int64  `json:"uploaded"`
	EstimatedTime   int    `json:"estimated_time"`
	PeersGetting    int    `json:"peers_getting_from_us"`
	Extract         bool   `json:"extract"`
	CurrentRatio    string `json:"current_ratio"`
	Size            int64  `json:"size"`
	UpSpeed         int64  `json:"up_speed"`
	Id              int64  `json:"id"`
	Source          string `json:"source"`
	Subscription_id int64  `json:"subscription_id"`
	StatusMessage   string `json:"status_message"`
	Status          string `json:"status"`
	DownSpeed       int64  `json:"down_speed"`
	PeersConnected  int    `json:"peers_connected"`
	Downloaded      int64  `json:"downloaded"`
	FileId          int64  `json:"file_id"`
	PeersSending    int    `json:"peers_sending_to_us"`
	PercentDone     int    `json:"percent_done"`
	IsPrivate       bool   `json:"is_private"`
	TrackerMessage  string `json:"tracker_message"`
	Name            string `json:"name"`
	CreatedAt       string `json:"created_at"`
	ErrorMessage    string `json:"error_message"`
	SaveParentId    int64  `json:"save_parent_id"`
	CallbackUrl     string `json:"callback_url"`
}

type Transfers struct {
	Status    string
	Transfers []Transfer
	Transfer  Transfer
}

type Disk struct {
	Available int64
	Used      int64
	Size      int64
}

type UserInfo struct {
	Username string
	Mail     string
	Disk     Disk
}

type Settings struct {
	Routing               string `json: "routing"`
	HideItemsShared       string `json: "hide_items_shared"`
	DefaultDownloadFolder int    `json: "default_download_folder"`
	SSLEnabled            bool   `json: "ssl_enabled"`
	IsInvisible           bool   `json: "is_invisible"`
	ExtractionDefault     string `json: "extraction_default"`
}

type Account struct {
	Status   string
	Info     UserInfo
	Settings Settings
}

type Friend struct {
	Name string
}

type Friends struct {
	Status  string
	Friends []Friend
	Friend  Friend
}

type Putio struct {
	OauthToken string
}

func (p *Putio) GetReqBody(path string) (bodybytes []byte, err error) {
	return p.GetReqBodyParams(path, nil)
}

func (p *Putio) GetReqBodyParams(path string, params url.Values) (bodybytes []byte, err error) {
	var url string
	if params == nil {
		url = BaseUrl + path + oauthparam + p.OauthToken
	} else {
		params.Add("oauth_token", p.OauthToken)
		url = BaseUrl + path + "?" + params.Encode()
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodybytes, nil
}

func (p *Putio) PostFilesReq(path string, data url.Values) (files *Files, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &files); err != nil {
		return nil, string(bodybytes), err
	}
	return files, string(bodybytes), nil
}

func (p *Putio) GetFilesReq(path string) (files *Files, jsonstr string, err error) {
	return p.GetFilesReqParams(path, nil)
}

func (p *Putio) GetFilesReqParams(path string, params url.Values) (files *Files, jsonstr string, err error) {
	if params == nil {
		params = make(url.Values)
	}
	bodybytes, err := p.GetReqBodyParams(path, params)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &files); err != nil {
		return nil, string(bodybytes), err
	}
	return files, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#files-list
func (p *Putio) FilesList() (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/list")
}

func (p *Putio) FilesListDir(dirNo int64) (files *Files, jsonstr string, err error) {
	var vals url.Values
	vals = make(url.Values)
	vals.Add("parent_id", strconv.FormatInt(int64(dirNo), 10))
	return p.GetFilesReqParams("files/list", vals)
}

// https://api.put.io/v2/docs/#files-search
func (p *Putio) FilesSearch(query string, pageno string) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/search/" + query + "/page/" + string(pageno))
}

// https://api.put.io/v2/docs/#files-create-folder
func (p *Putio) FilesCreateFolder(name string, parent_id int64) (files *Files, jsonstr string, err error) {
	data := make(url.Values)
	data.Set("name", name)
	data.Set("parent_id", string(parent_id))
	return p.PostFilesReq("files/create-folder", data)
}

// https://api.put.io/v2/docs/#files-id
func (p *Putio) FilesId(id int64) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/" + strconv.FormatInt(int64(id), 10))
}

// https://api.put.io/v2/docs/#files-delete
func (p *Putio) FilesDelete(file_id int64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/delete", url.Values{"file_ids": {strconv.FormatInt(int64(file_id), 10)}})
}

// https://api.put.io/v2/docs/#files-rename
func (p *Putio) FilesRename(file_id int64, name string) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/rename", url.Values{"file_id": {strconv.FormatInt(int64(file_id), 10)}, "name": {name}})
}

// https://api.put.io/v2/docs/#files-move
func (p *Putio) FilesMove(file_id int64, parent_id int64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/move", url.Values{"file_id": {strconv.FormatInt(int64(file_id), 10)}, "parent_id": {strconv.FormatInt(int64(parent_id), 10)}})
}

// https://api.put.io/v2/docs/#files-mp4-post
func (p *Putio) FilesMP4(id int64) (files *Files, jsonstr string, err error) {
	return p.PostFilesReq("files/"+strconv.FormatInt(int64(id), 10)+"/mp4", url.Values{"id": {strconv.FormatInt(int64(id), 10)}})
}

// https://api.put.io/v2/docs/#files-mp4-post
func (p *Putio) FilesMP4Status(id int64) (files *Files, jsonstr string, err error) {
	return p.GetFilesReq("files/" + strconv.FormatInt(int64(id), 10) + "/mp4")
}

// https://api.put.io/v2/docs/#files-id-download
// in this case we will just return the url to download from and leave it up to
// the client to actually download it. It's a redirect so can't use the usual request method
func (p *Putio) FilesDownload(id int64) (urlstr string, err error) {
	path := "download"

	url := BaseUrl + "files/" + strconv.FormatInt(int64(id), 10) + "/" + path + oauthparam + p.OauthToken
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	finalURL := resp.Request.URL.String()
	return finalURL, nil
}

func (p *Putio) PostTransfersReq(path string, data url.Values) (transfers *Transfers, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &transfers); err != nil {
		return nil, string(bodybytes), err
	}
	return transfers, string(bodybytes), nil
}

func (p *Putio) GetTransfersReq(path string) (transfers *Transfers, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &transfers); err != nil {
		return nil, string(bodybytes), err
	}
	return transfers, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#transfers-list
func (p *Putio) TransfersList() (transfers *Transfers, jsonstr string, err error) {
	return p.GetTransfersReq("transfers/list")
}

// https://api.put.io/v2/docs/#transfers-add
func (p *Putio) TransfersAdd(transfer_url string, save_parent_id int64, extract bool) (transfers *Transfers, jsonstr string, err error) {
	return p.PostTransfersReq("transfers/add", url.Values{"url": {transfer_url}, "save_parent_id": {strconv.FormatInt(int64(save_parent_id), 10)}, "extract": {strconv.FormatBool(extract)}})
}

// https://api.put.io/v2/docs/#transfers-add
func (p *Putio) TransfersCancel(transfer_id int64) (transfers *Transfers, jsonstr string, err error) {
	return p.PostTransfersReq("transfers/cancel", url.Values{"transfer_ids": {strconv.FormatInt(int64(transfer_id), 10)}})
}

// https://api.put.io/v2/docs/#transfers-id
func (p *Putio) TransfersId(id int64) (transfers *Transfers, jsonstr string, err error) {
	return p.GetTransfersReq("transfers/" + strconv.FormatInt(int64(id), 10))
}

func (p *Putio) GetAccountReq(path string) (account *Account, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &account); err != nil {
		return nil, string(bodybytes), err
	}
	return account, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#account-info
func (p *Putio) AccountInfo() (account *Account, jsonstr string, err error) {
	return p.GetAccountReq("account/info")
}

// https://api.put.io/v2/docs/#account-settings
func (p *Putio) AccountSettings() (account *Account, jsonstr string, err error) {
	return p.GetAccountReq("account/settings")
}

func (p *Putio) PostFriendsReq(path string, data url.Values) (friends *Friends, jsonstr string, err error) {
	posturl := BaseUrl + path + oauthparam + p.OauthToken
	resp, err := http.PostForm(posturl, data)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if err = json.Unmarshal(bodybytes, &friends); err != nil {
		return nil, string(bodybytes), err
	}
	return friends, string(bodybytes), nil
}

func (p *Putio) GetFriendReq(path string) (friends *Friends, jsonstr string, err error) {
	bodybytes, err := p.GetReqBody(path)
	if err != nil {
		return nil, string(bodybytes), err
	}
	if err = json.Unmarshal(bodybytes, &friends); err != nil {
		return nil, string(bodybytes), err
	}
	return friends, string(bodybytes), nil
}

// https://api.put.io/v2/docs/#friends-list
func (p *Putio) FriendsList() (friends *Friends, jsonstr string, err error) {
	return p.GetFriendReq("friends/list")
}

// https://api.put.io/v2/docs/#friends-username-request
func (p *Putio) FriendsRequest(username string) (friends *Friends, jsonstr string, err error) {
	return p.PostFriendsReq("friends/"+username+"/request", nil)
}

// https://api.put.io/v2/docs/#friends-username-deny
func (p *Putio) FriendsDeny(username string) (friends *Friends, jsonstr string, err error) {
	return p.PostFriendsReq("friends/"+username+"/deny", nil)
}

// https://api.put.io/v2/docs/#friends-waiting-requests
func (p *Putio) FriendsWaiting() (friends *Friends, jsonstr string, err error) {
	return p.GetFriendReq("friends/waiting-requests")
}

// NewPutio takes in the apps oauth information and gets the token that will be used for all other calls
// This function doesn't have to be used if you provied the OauthToken when creating a Putio struct.
func NewPutio(appid, appsecret, appredirect, usercode string) (*Putio, error) {
	// get the user token using the calling apps credentials
	url := "https://api.put.io/v2/oauth2/access_token?client_id=" + appid + "&client_secret=" + appsecret + "&grant_type=authorization_code&redirect_uri=" + appredirect + "&code=" + usercode
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// read in the body of the response
	defer resp.Body.Close()
	bodybytes, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// token returns as json result like { "access_token": "ABV9KDHN" }
	type oauthtoken struct {
		Access_token string
	}
	token := oauthtoken{}
	if err = json.Unmarshal(bodybytes, &token); err != nil {
		return nil, err
	}
	return &Putio{OauthToken: token.Access_token}, nil
}
