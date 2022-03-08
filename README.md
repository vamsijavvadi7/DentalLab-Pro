# testdb



type Details struct {
		Comp []*Result `json:"details"`
	}

	res := make([]*Result, 0)
	for rows.Next() {
		rt := new(Result)
		err := rows.Scan(&rt.CompetencyName, &rt.CompetencyId)

		if err != nil {
			panic(err)
		}
		res = append(res, rt)
	}
	p:=new(Details)
	p.Comp=make([]*Result, 0)
	for _, item := range res {
	p.Comp = append(p.Comp, &Result{CompetencyName :item.CompetencyName,CompetencyId: item.CompetencyId});
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(p)
