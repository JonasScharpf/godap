package main

import (
    "flag"
    "github.com/JonasScharpf/godap/godap"
    "log"
    "strings"
)

func main() {
    addr := flag.String("addr", ":389", "HTTP network address with port")
    flag.Parse()

    hs := make([]godap.LDAPRequestHandler, 0)

    // use a LDAPBindFuncHandler to provide a callback function to respond
    // to bind requests
    hs = append(hs, &godap.LDAPBindFuncHandler{LDAPBindFunc: func(binddn string, bindpw []byte) bool {
        if strings.Contains(binddn, "cn=user,") && string(bindpw) == "password" {
            return true
        }
        return false
    }})

    // use a LDAPSimpleSearchFuncHandler to reply to search queries
    hs = append(hs, &godap.LDAPSimpleSearchFuncHandler{LDAPSimpleSearchFunc: func(req *godap.LDAPSimpleSearchRequest) []*godap.LDAPSimpleSearchResultEntry {

        ret := make([]*godap.LDAPSimpleSearchResultEntry, 0, 1)

        // here we produce a single search result that matches whatever
        // they are searching for
        if req.FilterAttr == "uid" {
            ret = append(ret, &godap.LDAPSimpleSearchResultEntry{
                DN: "cn=" + req.FilterValue + "," + req.BaseDN,
                Attrs: map[string]interface{}{
                    "sn":            req.FilterValue,
                    "cn":            req.FilterValue,
                    "uid":           req.FilterValue,
                    "homeDirectory": "/home/" + req.FilterValue,
                    "objectClass": []string{
                        "top",
                        "posixAccount",
                        "inetOrgPerson",
                    },
                },
            })
        } else {
            log.Println("Currently nothing else than a '(uid=<name>)' query is handled")
        }

        return ret

    }})

    s := &godap.LDAPServer{
        Handlers: hs,
    }

    log.Println("Starting mock LDAP server on", *addr)
    err := s.ListenAndServe(*addr)

    if err != nil {
        log.Printf("Failed to start server. Error : %v", err)
    }
}
