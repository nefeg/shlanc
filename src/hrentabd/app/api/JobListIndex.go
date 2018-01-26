package api

type JobListIndex map[string]Job


func (jli *JobListIndex)Merge(m ...JobListIndex){

	for _, a := range m{
		for k, v := range a {
			(*jli)[k] = Job(v)
		}
	}
}
