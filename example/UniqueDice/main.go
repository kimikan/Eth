package main

import (
	"UniqueDice/ethutils"
	"UniqueDice/helpers"
	"UniqueDice/store"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("template/index.tmpl")
	if err != nil {
		writeInfoToClient(w, "failed", err.Error())
		return
	}
	w.Write(content)
}

func AmbrDetails(w http.ResponseWriter, req *http.Request) {
	/*
		ethaddr := req.PostFormValue("ethaddr")
		item, err := getMail(ethaddr)
		if err != nil {
			AmbrReturnResult(w, "用户不存在!")
		} else {
			content, err := ioutil.ReadFile("views/html/details.html")
			if err != nil {
				AmbrReturnResult(w, err.Error())
			}

			t, e := template.New("webpage").Parse(string(content))
			if e != nil {
				AmbrReturnResult(w, e.Error())
				return
			}

			err = t.Execute(w, item)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
		} */
}

func whitelistIndex(w http.ResponseWriter, r *http.Request) {
}

func checkWhitelist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		content, err := ioutil.ReadFile("template/whitelist/IsInWhitelist.tmpl")
		if err == nil {
			w.Write(content)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		tokenAddr := r.PostFormValue("tokenaddr")
		ok, e := ethutils.IsInUniqueDiceWhitelist(tokenAddr)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			result := "Token is in uniquedice whitelist"
			if !ok {
				result = "Token is NOT in uniquedice whitelist"
			}
			writeInfoToClient(w, "Success", result)
		}
	} //endif
}

func grantDelegate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		content, err := ioutil.ReadFile("template/grantdelegate.tmpl")
		if err == nil {
			w.Write(content)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		tokenAddr := r.PostFormValue("tokenaddr")
		tokenNum := r.PostFormValue("tokennumber")
		privKey := r.PostFormValue("privatekey")
		e := ethutils.GrantDelegate(privKey, tokenAddr, tokenNum)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			writeInfoToClient(w, "Success", "Granted to uniquedice")
		}
	} //endif
}

func playGame(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		content, err := ioutil.ReadFile("template/playgame.tmpl")
		if err == nil {
			w.Write(content)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
	} //endif
}

func writeInfoToClientByTemplate(w http.ResponseWriter, subject string, detail string, file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	type info struct {
		Subject string
		Detail  string
	}
	arg := &info{
		Subject: subject,
		Detail:  detail,
	}
	t, e := template.New("webpage").Parse(string(content))
	if e != nil {
		w.Write([]byte(e.Error()))
		return
	}

	err = t.Execute(w, arg)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func writeInfoToClient(w http.ResponseWriter, subject string, detail string) {
	writeInfoToClientByTemplate(w, subject, detail, "template/info.tmpl")
}

func about(w http.ResponseWriter, r *http.Request) {
	writeInfoToClientByTemplate(w, "About", "written by kimikan", "template/about.tmpl")
}

func getUniquediceTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		content, err := ioutil.ReadFile("template/gettokens.tmpl")
		if err == nil {
			w.Write(content)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		tokenAddr := r.PostFormValue("tokenaddr")
		num, e := ethutils.GetTokens(tokenAddr)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			writeInfoToClient(w, "Success", num.Text(10)+" tokens remain")
		}
	} //endif
}

func tokenList(w http.ResponseWriter, r *http.Request) {
	tokens, e := helpers.GetAllERC20Tokens2()
	if e != nil {
		writeInfoToClient(w, "Failed", e.Error())
		return
	}
	content, err := ioutil.ReadFile("template/tokenlist.tmpl")
	t, e := template.New("webpage").Parse(string(content))
	if e != nil {
		w.Write([]byte(e.Error()))
		return
	}

	err = t.Execute(w, tokens)
	if err != nil {
		writeInfoToClient(w, "Failed", err.Error())
	}
}

func whitelistModify(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		content, err := ioutil.ReadFile("template/whitelist/ModifyWhiteList.tmpl")
		if err == nil {
			w.Write(content)
		} else {
			w.Write([]byte(err.Error()))
		}
	} else {
		tokenAddr := r.PostFormValue("tokenaddr")
		privatekey := r.PostFormValue("privatekey")
		flag := r.PostFormValue("flag")
		//e := ethutils.ModifyWhiteList(tokenAddr, len(flag) <= 0)
		e := ethutils.ModifyWhiteListByPrivateKey(tokenAddr, privatekey, len(flag) <= 0)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			result := "Successfully added token into uniquedice whitelist!"
			if len(flag) > 0 {
				result = "Successfully removed token into uniquedice whitelist!"
			}
			writeInfoToClient(w, "Success", result)
		}
	} //endif
}

func requestWhitelist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//end if
		tokenAddr := r.PostFormValue("tokenaddr")
		detail := r.PostFormValue("detail")
		username := r.PostFormValue("username")
		email := r.PostFormValue("email")
		if !ethutils.IsAddress(tokenAddr) {
			writeInfoToClient(w, "Error", "Invalid token address!")
			return
		}
		if !strings.Contains(email, "@") {
			writeInfoToClient(w, "Error", "Invalid email")
			return
		}
		e := store.SubmitRequest(&store.Request{
			Username:     username,
			Email:        email,
			TokenAddress: tokenAddr,
			Description:  detail,
		})
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			writeInfoToClient(w, "Success", "Request successfully submitted!")
		}
	}
}

func approveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		reqs, e := store.GetOngoingRequests()
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
			return
		}
		content, err := ioutil.ReadFile("template/whitelist/requestlist.tmpl")
		t, e := template.New("webpage").Parse(string(content))
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		err = t.Execute(w, reqs)
		if err != nil {
			writeInfoToClient(w, "Failed", err.Error())
		}
	} else {
		writeInfoToClient(w, "Failed", "Not supported")
	}
}

func historyRequest(w http.ResponseWriter, r *http.Request) {
	reqs, e := store.GetApprovedRequests()
	if e != nil {
		writeInfoToClient(w, "Failed", e.Error())
		return
	}
	content, err := ioutil.ReadFile("template/whitelist/historyrequests.tmpl")
	t, e := template.New("webpage").Parse(string(content))
	if e != nil {
		w.Write([]byte(e.Error()))
		return
	}
	err = t.Execute(w, reqs)
	if err != nil {
		writeInfoToClient(w, "Failed", err.Error())
	}
}

const (
	kAdminPassword = "ambr"
)

func approveConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tokenaddress := r.URL.Query().Get("id")
		req, e := store.GetOngoingRequest(tokenaddress)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
			return
		}
		content, err := ioutil.ReadFile("template/whitelist/approveconfirm.tmpl")
		t, e := template.New("webpage").Parse(string(content))
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		err = t.Execute(w, req)
		if err != nil {
			writeInfoToClient(w, "Failed", err.Error())
		}
	} else {
		password := r.PostFormValue("password")
		if password != kAdminPassword {
			writeInfoToClient(w, "Error", "Invalid admin password")
			return
		}
		tokenaddress := r.PostFormValue("tokenaddr")
		e := store.ApproveRequest(tokenaddress)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			e2 := ethutils.ModifyWhiteList(tokenaddress, true)
			if e2 != nil {
				writeInfoToClient(w, "Failed", e2.Error())
			} else {
				writeInfoToClient(w, "Success", "Specific request approved!")
			}
		}
	}
}

func deleteConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tokenaddress := r.URL.Query().Get("id")
		req, e := store.GetOngoingRequest(tokenaddress)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
			return
		}
		content, err := ioutil.ReadFile("template/whitelist/deleteconfirm.tmpl")
		t, e := template.New("webpage").Parse(string(content))
		if e != nil {
			w.Write([]byte(e.Error()))
			return
		}
		err = t.Execute(w, req)
		if err != nil {
			writeInfoToClient(w, "Failed", err.Error())
		}
	} else {
		password := r.PostFormValue("password")
		if password != kAdminPassword {
			writeInfoToClient(w, "Error", "Invalid admin password")
			return
		}
		tokenaddress := r.PostFormValue("tokenaddr")
		e := store.DeleteOngoingRequest(tokenaddress)
		if e != nil {
			writeInfoToClient(w, "Failed", e.Error())
		} else {
			writeInfoToClient(w, "Success", "Specific request deleted!")
		}
	}
}

func setupWhitelistHandles() {
	http.HandleFunc("/whitelist/checkwhitelist", checkWhitelist)
	http.HandleFunc("/whitelist", whitelistIndex)
	http.HandleFunc("/whitelist/modify", whitelistModify)
	http.HandleFunc("/whitelist/request", requestWhitelist)
	http.HandleFunc("/whitelist/approverequest", approveRequest)
	http.HandleFunc("/whitelist/historyrequests", historyRequest)
	http.HandleFunc("/whitelist/request/approve", approveConfirm)
	http.HandleFunc("/whitelist/request/delete", deleteConfirm)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc("/register", AmbrRegister)
	/*http.HandleFunc("/details", AmbrDetails)
	http.HandleFunc("/check", AmbrCheck)
	http.HandleFunc("/download", AmbrDownload)

		p.Mux.HandleFunc("/zid/details", biz.ZidDetails)
		p.Mux.HandleFunc("/zid/check", biz.ZidCheck)
		p.Mux.HandleFunc("/zid/download", biz.ZidDownload)
		p.Mux.HandleFunc("/zid", biz.ZidIndex)
	*/
	setupWhitelistHandles()
	http.HandleFunc("/getuniquedicetokens", getUniquediceTokens)
	http.HandleFunc("/grant", grantDelegate)
	http.HandleFunc("/playgame", playGame)
	http.HandleFunc("/tokenlist", tokenList)
	http.HandleFunc("/about", about)
	http.HandleFunc("/", IndexPage)
	http.ListenAndServe(":8002", nil)
}
