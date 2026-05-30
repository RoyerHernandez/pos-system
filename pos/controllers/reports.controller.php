<?php

class ReportController{

	static public function ctrSalesByDateRange($startDate, $endDate){
		return ReportModel::mdlSalesByDateRange($startDate, $endDate);
	}

	static public function ctrSalesSummary($startDate, $endDate){
		return ReportModel::mdlSalesSummary($startDate, $endDate);
	}

	static public function ctrTopProducts($startDate, $endDate, $limit = 10){
		return ReportModel::mdlTopProducts($startDate, $endDate, $limit);
	}

	static public function ctrSalesByPaymentMethod($startDate, $endDate){
		return ReportModel::mdlSalesByPaymentMethod($startDate, $endDate);
	}

	static public function ctrCashRegisterReport($startDate, $endDate){
		return ReportModel::mdlCashRegisterReport($startDate, $endDate);
	}

	static public function ctrDailySalesBreakdown($startDate, $endDate){
		return ReportModel::mdlDailySalesBreakdown($startDate, $endDate);
	}

}
