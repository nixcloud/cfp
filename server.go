package main

/*
 * cookie setzen
 * form
 * websocket?
 * save submission as json
 * send submission via email
 */

import (
    "github.com/gocraft/web"
    "fmt"
    "net/http"
    "os"
    "text/template"
    "gopkg.in/gomail.v2"
    "strings"
    "time"
    "bytes"
    "io/ioutil"
    "crypto/md5"

//     "bufio"
//     "math/rand"

//     "github.com/dchest/captcha"
//     "io"
//     "golang.org/x/net/websocket"
)

type Context struct {
    HelloCount int
}

// var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
// 
// // generates a random string of fixed size
// func srand(size int) string {
//     buf := make([]byte, size)
//     for i := 0; i < size; i++ {
//         buf[i] = alpha[rand.Intn(len(alpha))]
//     }
//     return string(buf)
// }

func (c *Context) Submit(rw web.ResponseWriter, req *web.Request) {
  
   noItems := struct {
     AGB                    string
     Name                   string
     Email                  string
     Datum                  string
     V25                    string
     V55                    string
     VortragAlsWS           string
     LT5                    string
     LT10                   string
     Workshop               string
     MitschnittZugestimmt   string
     Titel                  string
     Buzzword               string
     Inhalt                 string
     Vorwissen              string
     Weblinks               string
     Sonstiges              string
    }{};
  
    req.ParseForm()  //Parse url parameters passed, then parse the response packet for the POST body (request body)

    ti := time.Now();
    noItems.Datum = ti.UTC().Format(time.UnixDate);
    for key, value := range req.Form {
      zz := strings.Replace(strings.Join(value, " "),"\"", "\\\"", -1);
//       zz := strings.Join(value, " ");
      switch {
        case key == "mitschnittZugestimmt":
          noItems.MitschnittZugestimmt = zz;
        case key == "type.v25":
          noItems.V25 = zz
        case key == "type.v55":
          noItems.V55 = zz
        case key == "type.vortrag-als-ws":
          noItems.VortragAlsWS = zz
        case key == "type.lt10":
          noItems.LT10 = zz
        case key == "type.lt5":
          noItems.LT5 = zz
        case key == "type.workshop":
          noItems.Workshop = zz
        case key == "titel":
          noItems.Titel = zz
        case key == "buzzword":
          noItems.Buzzword = zz
        case key == "vorwissen":
          noItems.Vorwissen = zz
        case key == "AGBakzeptiert":
          noItems.AGB = zz
        case key == "email":
          noItems.Email = zz
        case key == "inhalt":
          noItems.Inhalt = zz
        case key == "weblinks":
          noItems.Weblinks = zz
        case key == "name":
          noItems.Name = zz
        case key == "sonstiges":
          noItems.Sonstiges = zz
      }
      
//         fmt.Fprint(rw,"key: ", key, ", value: ", value,"\n")
//         fmt.Print (rw,"key: ", key, ", value: ", value,"\n")
    }

const jsonExportTemplate = 
`
{
  "AGB": "{{.AGB}}",
  "name": "{{.Name}}",
  "email": "{{.Email}}",
  "einreichungAm": "{{.Datum}}",
  "type": {
    "v25": {{.V25}},
    "v55": {{.V55}},
    "alsWSmoeglich": {{.VortragAlsWS}},
    "lightning5": {{.LT5}},
    "workshop": {{.Workshop}}
   },
  "mitschnittZugestimmt": {{.MitschnittZugestimmt}},
  "titel": "{{.Titel}}",
  "buzzword": "{{.Buzzword}}",
  "inhalt": "{{.Inhalt}}",
  "vorwissen": "{{.Vorwissen}}",
  "weblinks": "{{.Weblinks}}",
  "sonstiges": "{{.Sonstiges}}"
}
`

// Create a new template and parse the letter into it.
  t := template.Must(template.New("jsonExportTemplate").Parse(jsonExportTemplate))


  buff := bytes.NewBufferString("")
  err := t.Execute(buff, noItems)
  if err != nil {
    fmt.Println("executing template:", err)
  }
  
//   fmt.Println(string(buff.Bytes()))

    // save data to disk
    var a = fmt.Sprintf("%x", md5.Sum(buff.Bytes()))

    _ = ioutil.WriteFile("incomming/" + a, buff.Bytes(), 0644)
    
//     https://godoc.org/gopkg.in/gomail.v2
    m := gomail.NewMessage()
    m.SetHeader("From", "example@example.com")
    m.SetHeader("To", "cfp-signup-2016@tuebix.org")
//     m.SetHeader("CC", "js@lastlog.de")
    m.SetHeader("Subject", "tuebix16 beitrag: " + noItems.Name + ": " + noItems.Titel)
    
    p := `Hallo,
    
    ich bin der server nixcloud.io und sende neue tuebix anfragen:
    
    ================= <json> =================
    `
    
    p += string(buff.Bytes())    
    
    p += `
    
    ================= </json> =================
    
    falls es probleme mit diesem dienst gibt, bitte an:
    joachim schiele <js@lastlog.de> wenden.
    
    glg,
    joachim
    `
    m.SetBody("text/plain", p)



    
    d := gomail.NewPlainDialer("example.com", 587, "example@example.com", "examplepassword")
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}

func (c *Context) Main(rw web.ResponseWriter, req *web.Request) {
  t, err := template.New("index.tmpl").
  ParseFiles("index.tmpl")
  if err != nil {
    fmt.Println(err)
    panic(err)
    os.Exit(1)
  }
  noItems := struct {
//     Title       string
  }{
//     Title:        title,
  }
  err = t.Execute(rw, noItems)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}

func main() {
    currentRoot, _ := os.Getwd() 
    currentRoot +=  "/public"

    router := web.New(Context{}).               // Create your router
    Middleware(web.LoggerMiddleware).           // Use some included middleware
    Middleware(web.ShowErrorsMiddleware).       // ...
    Get("/", (*Context).Main).                  // Add a route
    Post("/submit", (*Context).Submit).
    Middleware(web.StaticMiddleware(currentRoot, web.StaticOption{Prefix: "/public"}))

    http.ListenAndServe("localhost:3000", router)   // Start the server!
}
