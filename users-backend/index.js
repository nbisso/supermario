var express = require("express");

var app = express();

var os = require("os");
const mongoose = require('mongoose');

let mongoURL = 'mongodb://mongo-0.mongo:27017,mongo-1.mongo:27017,mongo-2.mongo:27017/auth'

let Example = require('./model')


var morgan = require('morgan');
var cors = require('cors')

app.use(cors());

app.use(morgan("dev"));

app.get("/users", (req, res) => {
    res.json([
        {
            id: 1,
            name: "pedro",
            lastname: "gomez",
            age: 32
        },
        {
            id: 2,
            name: "juan",
            lastname: "gomez",
            age: 32
        }
    ])
})





app.get("/users/example", (req, res) => {
    let userData = {
        id: Math.floor(Math.random() * 100000),
        data: "daslkndalsk",
    }
    new Example(userData).save()
        .then(item => {
            res.send("item saved to database");
        })
        .catch(err => {
            console.log(err)
            res.status(400).send(err);
        });
})

app.get("/users/examples", (req, res) => {

    Example.find({}, function (err, result) {
        if (err) {
            res.send(err);
        } else {
            res.send(result);
        }
    });
})

app.get("*", (req, res) => {
    var hostname = os.hostname();
    res.header("X-User", req.header("x-user"))
    res.send("Soy el host:" + hostname + "y el user: " + req.header("x-user"))
})


mongoose.connect(mongoURL, { useNewUrlParser: true })
    .then(
        () => {
            console.log("connected to mongo");
        }
    ).catch((error) => {
        console.log("unable to connect to mongoDB")
        console.log(mongoURL);
    });



app.listen(8080, () => {
    console.log("server init on localhost:8080")
})