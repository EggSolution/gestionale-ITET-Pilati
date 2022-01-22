var express = require('express');
var router = express.Router();
var db1 = require('../database/db1');

// indirizzamento in caoso di nessuna pagina
router.get('/', function(req, res) {
  	res.render('index', {title: 'Home'});
});

// indirizzamento in caso di /login
router.get('/login', function(req, res) {
	res.render('login', {title: 'Login'});
});

router.get("/register", function(req, res){
	res.render('register', {title: 'Register'});
});

router.post('/dashboard', function(req, res){

	const email = req.body.email;
	const password = req.body.password;

	var query = 'SELECT * FROM user WHERE email="' + email + '" AND password="' + password + '"';
	db1.query(query, function(err, result, fields){
		if(err){
			console.log(err);
		} else {
			if(result == ""){
				res.render('login', {title: "Login error", errore: true});
			} else {
				console.log(result[0].name)
				res.render('dashboard', {title: "Dashboard", nomeUtente: result[0].name});
			}
		}
	});
});

router.post('/register', function(req, res){

	const nome = req.body.nome;
	const email = req.body.email;
	const password1 = req.body.pass1;
	const password2 = req.body.pass2;
	var enUsata = false;

	/*
	db1.query("SELECT * FROM user WHERE email='"+email+"' OR name='"+nome+"'", function(err, result){
		if(err){
			console.log(err);
		} else {
			if(result[0].id != undefined){
				enUsata = true;
			} else {
				enUsata = false;
			}
		}
	});
	*/

	// da sistemare gli errori non arrivano in modo ordinato
	if (password1 === password2){
		if(enUsata == false){
			if(nome != "" || email != "" || password1 != "" || password2 != ""){
				var data = new Date();
				var dataFormattata = data.getDay()+"/"+data.getMonth()+"/"+data.getFullYear();
					
				var query = "INSERT INTO user (name, privileges, data, password, email) VALUE ('"+nome+"',"+1+",'"+dataFormattata+"','"+password2+"','"+email+"')";
				db1.query(query, function(err, result, fields){
					if(err){
						console.log(err);
					} else {
						console.log(nome + " si è loggato");
						res.redirect('/login');
					}
				});
			} else {
				res.render('register', {title: "Register error", errore: "errore: riempire tutti i campi"});
			}
		} else if(enUsata == true) {
			res.render('register', {title: "Register error", errore: "errore: nome utente o password già utilizati"});
		}
	} else {
		res.render('register', {title: "Register error", errore: "errore: password diverse"});
	}

});

module.exports = router;
