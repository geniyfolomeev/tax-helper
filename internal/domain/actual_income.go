package domain

import "time"

type ActualIncome struct {
	PrevMonth []*Income
	CurMonth  []*Income
}

func NewActualIncome(incomes []*Income) *ActualIncome {
	prevFrom, prevTo, curFrom, curTo := GetActualTimeIntervals()
	prevIncomes := make([]*Income, 0)
	curIncomes := make([]*Income, 0)

	for _, income := range incomes {
		switch {
		case !income.Date.Before(prevFrom) && !income.Date.After(prevTo):
			prevIncomes = append(prevIncomes, income)
		case !income.Date.Before(curFrom) && !income.Date.After(curTo):
			curIncomes = append(curIncomes, income)
		}
	}

	return &ActualIncome{
		PrevMonth: prevIncomes,
		CurMonth:  curIncomes,
	}
}

func GetActualTimeIntervals() (prevFrom, prevTo, curFrom, curTo time.Time) {
	now := currentTimeFn()

	curFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	curTo = curFrom.AddDate(0, 1, 0).Add(-time.Nanosecond)

	prevFrom = curFrom.AddDate(0, -1, 0)
	prevTo = curFrom.Add(-time.Nanosecond)
	return
}
