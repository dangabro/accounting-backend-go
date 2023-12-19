package dao

const starterSqlTrialBalance = `select t.transaction_id,
			   t.transaction_date,
			   tp.account_id,
			   tp.amount_debit,
			   tp.amount_credit
		from transaction t
				 left outer join transaction_position tp on t.transaction_id = tp.transaction_id`
