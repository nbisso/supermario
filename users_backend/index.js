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
    res.send("Soy el host:" + hostname)
})

app.listen(8080, () => {
    console.log("server init on localhost:8080")
})