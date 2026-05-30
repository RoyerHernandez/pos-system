<?php

class Connection{

	static public function connect(){

		$host = getenv('DB_HOST') ?: 'localhost';
		$dbname = getenv('DB_NAME') ?: 'pos';
		$user = getenv('DB_USER') ?: 'root';
		$pass = getenv('DB_PASS') ?: '';

		$link = new PDO(
			"mysql:host=$host;dbname=$dbname",
			$user,
			$pass,
			array(PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION)
		);

		$link->exec("set names utf8");

		return $link;

	}

}