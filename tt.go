package main

import (
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "os"
  // "log"
  // "encoding/json"
  "flag"
  "time"
  "log"
)

type Time struct {
  Id              bson.ObjectId             `json:"id" bson:"_id,omitempty"`
  Name            string                    `json:"name" bson:"name"`
  Start           time.Time                 `json:"start" bson:"start"`
  End             time.Time                 `json:"end" bson:"end,omitempty"`
}

func main() { 
  sess, err := mgo.Dial("127.0.0.1")
  if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
  }
  defer sess.Close()
 
  sess.SetSafe(&mgo.Safe{})

  c := sess.DB("test").C("time")

  flag.Parse()
  switch verb := flag.Arg(0); verb {
  case "start":
    if len(flag.Args()) == 2 {
      var count int
      name := flag.Arg(1)
      query := bson.M{"name": name, "start": bson.M{"$exists": true}, "end": bson.M{"$exists": false}}
      count, err = c.Find(query).Count()
      if err != nil {
        fmt.Printf("got an error finding a doc %v\n")
        os.Exit(1)
      }
      if count == 0 {
        p := new(Time)
        p.Start = time.Now()
        p.Name = name
        err = c.Insert(p)
        if err != nil {
          log.Fatal(err)
        } 
        fmt.Println(p.Name)
        fmt.Println(p.Start)
      } else {
        fmt.Println("Você já possui uma atividade aberta.")
      }
    } else {
      fmt.Println("Passe um parametro")
    }
  case "end":
    if len(flag.Args()) == 2 {
      var count int
      name := flag.Arg(1)
      query := bson.M{"name": name, "start": bson.M{"$exists": true}, "end": bson.M{"$exists": false}}

      count, err = c.Find(query).Count()
      if count == 0 {
        fmt.Println("Não existe atividade aberta com esse nome.")
      } else {
        update := bson.M{"$set": bson.M{"end": time.Now()}}
        err = c.Update(query, update)
        if err != nil {
          os.Exit(1)
        }
      }
    } else {
      fmt.Println("error")
    }
  default:
    if len(flag.Args()) == 1 {
      name := flag.Arg(0)
      query := bson.M{"name": name}
      var result []Time
      err = c.Find(query).All(&result)
      if err != nil {
        fmt.Printf("got an error finding a doc %v\n")
        os.Exit(1)
      }
      if len(result) == 0 {
        fmt.Println("Você não possui tempos para este nome.")
      } else {
        var total int64
        fmt.Printf("Start\t\tEnd\t\tTotal\n")
        for _, t := range result {
          var duration time.Duration
          duration = t.End.UTC().Sub(t.Start.UTC())
          total += duration.Nanoseconds()
          fmt.Printf("%v\t%v\t%v\n", t.Start.Format("15:04:05"), t.End.Format("15:04:05"), duration)
        }
        totalDuration := time.Duration(total)
        fmt.Printf("\t\t\t\t%v\n", totalDuration)
      }

    } else {
      fmt.Println("Passe um parametro")
    }
  }
}