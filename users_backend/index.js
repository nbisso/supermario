var express = require("express");

var app = express();

var os = require("os");


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

app.get("*", (req, res) => {
    var hostname = os.hostname();
    res.header("X-User", req.header("x-user"))
    res.send("Soy el host:" + hostname + "y el user: " + req.header("x-user"))
})

app.listen(8080, () => {
    console.log("server init on localhost:8080")
})