package api

type JobListTime map[int64]JobListIndex

func (jlt *JobListTime)ToIndexList() JobListIndex{

	indexList := JobListIndex{}

	for _, li := range *jlt{

		indexList.Merge(li)
	}

	return indexList
}