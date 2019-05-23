package main

import "github.com/isaac/app"

func main() {
	var app app.App
	app.Initialize()
	app.InitializeRoutes()
	app.Run("9000")
	app.Close()
}
/*
INSERT INTO applicant(id_number, surname, other_name, no_of_dependents, mobile_number, alternative_number) VALUES(360827132, 'Henry', 'Kariuki', 3, '0707074707' '070530707007');
*/
