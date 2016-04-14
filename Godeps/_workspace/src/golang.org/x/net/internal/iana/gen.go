// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

//go:generate go run gen.go

// This program generates internet protocol constants and tables by
// reading IANA protocol registries.
package main

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "go/format"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"
)

var registries = []struct {
    url   string
    parse func(io.Writer, io.Reader) error
}{
    {
        "http://www.iana.org/assignments/dscp-registry/dscp-registry.xml",
        parseDSCPRegistry,
    },
    {
        "http://www.iana.org/assignments/ipv4-tos-byte/ipv4-tos-byte.xml",
        parseTOSTCByte,
    },
    {
        "http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xml",
        parseProtocolNumbers,
    },
}

func main() {
    var bb bytes.Buffer
    fmt.Fprintf(&bb, "// go generate gen.go\n")
    fmt.Fprintf(&bb, "// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT\n\n")
    fmt.Fprintf(&bb, "// Package iana provides protocol number resources managed by the Internet Assigned Numbers Authority (IANA).\n")
    fmt.Fprintf(&bb, `package iana // import "golang.org/x/net/internal/iana"` + "\n\n")
    for _, r := range registries {
        resp, err := http.Get(r.url)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            fmt.Fprintf(os.Stderr, "got HTTP status code %v for %v\n", resp.StatusCode, r.url)
            os.Exit(1)
        }
        if err := r.parse(&bb, resp.Body); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        fmt.Fprintf(&bb, "\n")
    }
    b, err := format.Source(bb.Bytes())
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    if err := ioutil.WriteFile("const.go", b, 0644); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func parseDSCPRegistry(w io.Writer, r io.Reader) error {
    dec := xml.NewDecoder(r)
    var dr dscpRegistry
    if err := dec.Decode(&dr); err != nil {
        return err
    }
    drs := dr.escape()
    fmt.Fprintf(w, "// %s, Updated: %s\n", dr.Title, dr.Updated)
    fmt.Fprintf(w, "const (\n")
    for _, dr := range drs {
        fmt.Fprintf(w, "DiffServ%s = %#x", dr.Name, dr.Value)
        fmt.Fprintf(w, "// %s\n", dr.OrigName)
    }
    fmt.Fprintf(w, ")\n")
    return nil
}

type dscpRegistry struct {
    XMLName     xml.Name `xml:"registry"`
    Title       string   `xml:"title"`
    Updated     string   `xml:"updated"`
    Note        string   `xml:"note"`
    RegTitle    string   `xml:"registry>title"`
    PoolRecords []struct {
        Name  string `xml:"name"`
        Space string `xml:"space"`
    } `xml:"registry>record"`
    Records     []struct {
        Name  string `xml:"name"`
        Space string `xml:"space"`
    } `xml:"registry>registry>record"`
}

type canonDSCPRecord struct {
    OrigName string
    Name     string
    Value    int
}

func (drr *dscpRegistry) escape() []canonDSCPRecord {
    drs := make([]canonDSCPRecord, len(drr.Records))
    sr := strings.NewReplacer(
        "+", "",
        "-", "",
        "/", "",
        ".", "",
        " ", "",
    )
    for i, dr := range drr.Records {
        s := strings.TrimSpace(dr.Name)
        drs[i].OrigName = s
        drs[i].Name = sr.Replace(s)
        n, err := strconv.ParseUint(dr.Space, 2, 8)
        if err != nil {
            continue
        }
        drs[i].Value = int(n) << 2
    }
    return drs
}

func parseTOSTCByte(w io.Writer, r io.Reader) error {
    dec := xml.NewDecoder(r)
    var ttb tosTCByte
    if err := dec.Decode(&ttb); err != nil {
        return err
    }
    trs := ttb.escape()
    fmt.Fprintf(w, "// %s, Updated: %s\n", ttb.Title, ttb.Updated)
    fmt.Fprintf(w, "const (\n")
    for _, tr := range trs {
        fmt.Fprintf(w, "%s = %#x", tr.Keyword, tr.Value)
        fmt.Fprintf(w, "// %s\n", tr.OrigKeyword)
    }
    fmt.Fprintf(w, ")\n")
    return nil
}

type tosTCByte struct {
    XMLName  xml.Name `xml:"registry"`
    Title    string   `xml:"title"`
    Updated  string   `xml:"updated"`
    Note     string   `xml:"note"`
    RegTitle string   `xml:"registry>title"`
    Records  []struct {
        Binary  string `xml:"binary"`
        Keyword string `xml:"keyword"`
    } `xml:"registry>record"`
}

type canonTOSTCByteRecord struct {
    OrigKeyword string
    Keyword     string
    Value       int
}

func (ttb *tosTCByte) escape() []canonTOSTCByteRecord {
    trs := make([]canonTOSTCByteRecord, len(ttb.Records))
    sr := strings.NewReplacer(
        "Capable", "",
        "(", "",
        ")", "",
        "+", "",
        "-", "",
        "/", "",
        ".", "",
        " ", "",
    )
    for i, tr := range ttb.Records {
        s := strings.TrimSpace(tr.Keyword)
        trs[i].OrigKeyword = s
        ss := strings.Split(s, " ")
        if len(ss) > 1 {
            trs[i].Keyword = strings.Join(ss[1:], " ")
        } else {
            trs[i].Keyword = ss[0]
        }
        trs[i].Keyword = sr.Replace(trs[i].Keyword)
        n, err := strconv.ParseUint(tr.Binary, 2, 8)
        if err != nil {
            continue
        }
        trs[i].Value = int(n)
    }
    return trs
}

func parseProtocolNumbers(w io.Writer, r io.Reader) error {
    dec := xml.NewDecoder(r)
    var pn protocolNumbers
    if err := dec.Decode(&pn); err != nil {
        return err
    }
    prs := pn.escape()
    prs = append([]canonProtocolRecord{{
        Name:  "IP",
        Descr: "IPv4 encapsulation, pseudo protocol number",
        Value: 0,
    }}, prs...)
    fmt.Fprintf(w, "// %s, Updated: %s\n", pn.Title, pn.Updated)
    fmt.Fprintf(w, "const (\n")
    for _, pr := range prs {
        if pr.Name == "" {
            continue
        }
        fmt.Fprintf(w, "Protocol%s = %d", pr.Name, pr.Value)
        s := pr.Descr
        if s == "" {
            s = pr.OrigName
        }
        fmt.Fprintf(w, "// %s\n", s)
    }
    fmt.Fprintf(w, ")\n")
    return nil
}

type protocolNumbers struct {
    XMLName  xml.Name `xml:"registry"`
    Title    string   `xml:"title"`
    Updated  string   `xml:"updated"`
    RegTitle string   `xml:"registry>title"`
    Note     string   `xml:"registry>note"`
    Records  []struct {
        Value string `xml:"value"`
        Name  string `xml:"name"`
        Descr string `xml:"description"`
    } `xml:"registry>record"`
}

type canonProtocolRecord struct {
    OrigName string
    Name     string
    Descr    string
    Value    int
}

func (pn *protocolNumbers) escape() []canonProtocolRecord {
    prs := make([]canonProtocolRecord, len(pn.Records))
    sr := strings.NewReplacer(
        "-in-", "in",
        "-within-", "within",
        "-over-", "over",
        "+", "P",
        "-", "",
        "/", "",
        ".", "",
        " ", "",
    )
    for i, pr := range pn.Records {
        if strings.Contains(pr.Name, "Deprecated") ||
        strings.Contains(pr.Name, "deprecated") {
            continue
        }
        prs[i].OrigName = pr.Name
        s := strings.TrimSpace(pr.Name)
        switch pr.Name {
        case "ISIS over IPv4":
            prs[i].Name = "ISIS"
        case "manet":
            prs[i].Name = "MANET"
        default:
            prs[i].Name = sr.Replace(s)
        }
        ss := strings.Split(pr.Descr, "\n")
        for i := range ss {
            ss[i] = strings.TrimSpace(ss[i])
        }
        if len(ss) > 1 {
            prs[i].Descr = strings.Join(ss, " ")
        } else {
            prs[i].Descr = ss[0]
        }
        prs[i].Value, _ = strconv.Atoi(pr.Value)
    }
    return prs
}
