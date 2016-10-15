//Lets require/import the HTTP module
var http = require('http');
var dispatcher = require('httpdispatcher');

//Lets define a port we want to listen to
const PORT=8081; 


//Lets use our dispatcher
function handleRequest(request, response){
    try {
        //log the request on console
        console.log(request.url);
        //Disptach
        dispatcher.dispatch(request, response);
    } catch(err) {
        console.log(err);
    }
}


var lolcats = [ "https://cutecatshq.com/wp-content/uploads/2016/03/This-beauty-is-strolling-around-my-local-cat-cafe.jpg",
"http://farm5.static.flickr.com/4138/4746437053_3c373a33ca.jpg",
"https://bw-2e2c4bf7ceaa4712a72dd5ee136dc9a8-bwcore.s3.amazonaws.com/articles/full_14229.jpg"
]

//A sample GET request    
dispatcher.onGet("/random_lolcat", function(req, res) {
    res.writeHead(200, {'Content-Type': 'text/plain'});
    var id = Math.floor((Math.random() * lolcats.length));
    res.end(lolcats[id]);
});    

//Create a server
var server = http.createServer(handleRequest);

//Lets start our server
server.listen(PORT, function(){
    //Callback triggered when server is successfully listening. Hurray!
    console.log("Server listening on: http://localhost:%s", PORT);
});