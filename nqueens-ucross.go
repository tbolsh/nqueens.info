package main
import(
	//"sync"
	"sort"
	"time"
	"flag"
	"fmt"
	"io"
	//"os"
)

var nq = flag.Int("n", 8, "N - the size of the board and amounts of queen to place. N > 0!")
var w  = flag.Bool("w", false, "Use modifyed (recursive->iterative) N. Wirth's algorithm, otherwise use T. Bolshakov's algorithm")
func main() {
	flag.Parse()
	fmt.Printf("Trying to solve N Queens problem for %d Queens on %d x %d board\n", *nq, *nq, *nq);  
	if *nq <= 0 {
		fmt.Println("N=%d, %d<=0, problem solved: no chessboard, no queens!\n", *nq)
		return
	}

	var solutions []Solution
	start := time.Now()
	if *w {
		fmt.Println("Using modifyed (recursive->iterative) N. Wirth algorithm.");	
		solutions = Wirth(*nq)
	} else {
		fmt.Println("Using T. Bolshakov's parallel algorithm.")
		solutions = Bolshakov(*nq)
	}
	dur := time.Since(start)
	if len(solutions) > 0 {
		fmt.Printf("All (%d) solutions for %d took %s\n", len(solutions), *nq, dur.String())
		//printSolutions(os.Stdout, solutions)
	}else{
		fmt.Printf("Solution for %d was not found!\n", *nq)
	}
	/*
	rowsb, rows := []int{1,2,3,4}, []int{0, 0, 0, 0}
	fmt.Println("")
	for b := 0; b < 3; b++ {
		for a := 0; a < 4; a++ {
			for r, c := range rowsb { rows[(r+a)%4] = c }
			for _, c := range rows {
				fmt.Printf("%d, ", c)
			}
			fmt.Println("")
		}
		rowsb[b], rowsb[(b+1)%4] = rowsb[(b+1)%4], rowsb[b] 
	}
	*/
}

type Solution struct {
	N int
	Rows []int
}

type SolutionW struct {
	Solution
	Cols, Diags, RDiags []bool
}

func compare(s1, s2 *Solution) bool{
	if(s1.N == s2.N){
		for i,s := range s1.Rows {
			if s!=s2.Rows[i] { return s < s2.Rows[i] } 
		}
	}
	return s1.N < s2.N
}

type solutionsSorter struct {
	sols []Solution
	comp func(s1, s2 *Solution) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *solutionsSorter) Len() int {
	return len(s.sols)
}
// Swap is part of sort.Interface.
func (s *solutionsSorter) Swap(i, j int) {
	s.sols[i], s.sols[j] = s.sols[j], s.sols[i]
}
// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *solutionsSorter) Less(i, j int) bool {
	return s.comp(&s.sols[i], &s.sols[j])
}


func equals(x, y []int) bool{
	if len(x)!=len(y) {
		return false
	}
	for i, e := range x {
		if e!=y[i] {
			return false
		}
	}
	return true
}

func addAllSolutions(s Solution, all []Solution) []Solution{
	all = addSolution(s, all)
	rows := s.Rows
	s.Rows = make([]int, s.N, s.N); 
	for r,c := range rows { s.Rows[c] = r }
	all = addSolution(s, all)
	for r,c := range rows { s.Rows[s.N-1-r] = c }
	all = addSolution(s, all)
	for r,c := range rows { s.Rows[s.N-1-c] = r }
	all = addSolution(s, all)
	return all
}

func addSolution(s2add Solution, all []Solution) []Solution{
	for _, s := range all {
		if equals( s2add.Rows, s.Rows ){
			return all
		}
	}
	var s1 Solution
	s1.N = s2add.N
	s1.Rows = make([]int, s1.N, s1.N); 			
	copy(s1.Rows, s2add.Rows)
	return append(all, s1)
}

func registerSolution(n int, rows []int/*, cols, diags, rdiags []bool*/) Solution{
	var s Solution
	s.N = n
	s.Rows   = make([]int,  n,     n); 			copy(s.Rows,   rows)
	return s
}

func newSolution(n int) SolutionW{
	var s SolutionW
	s.N = n
	s.Rows   = make([]int,  n,     n)
	s.Cols   = make([]bool, n,     n)
	s.Diags  = make([]bool, 2*n-1, 2*n-1)
	s.RDiags = make([]bool, 2*n-1, 2*n-1)
	reset( &s )
	return s
}

func reset(s *SolutionW){
	for i := range s.Rows {
		s.Rows[i], s.Cols[i] = -1, true
	}
	for i := range s.Diags {
		s.RDiags[i], s.Diags[i] = true, true
	}
}

func printSolution(w io.Writer, s Solution){
	for i,r := range s.Rows {
		if i == 0  { 
			fmt.Fprintf(w, "\t%d, ", r+1)
		} else if i == s.N-1 {
			fmt.Fprintf(w, "%d\n", r+1)
		} else {
			fmt.Fprintf(w, "%d, ", r+1)
		}
	}
}

func printSolutions(w io.Writer, solutions []Solution){
	if len(solutions)>0 {
		fmt.Fprintf(w, "%d (%d):\n", solutions[0].N, len(solutions))
		for _,s := range solutions {
			printSolution(w, s)
		}
	}else{
		fmt.Fprintln(w, "No solutions!\n")
	}
}

// Наивное, прямое, но не рекурсивное решение, скопировано из Вирта.
func Wirth(n int) []Solution{ 
	rows   := make([]int,  n,     n)
	cols   := make([]bool, n,     n)
	diags  := make([]bool, 2*n-1, 2*n-1)
	rdiags := make([]bool, 2*n-1, 2*n-1)
	solved := false
	solutions := make([]Solution, 0, 1000) // binomial coefficient 2*n-1, n ... 
	for i:=range rows {
		rows[i], cols[i] = 0, true
	}
	for i:=range diags {
		diags[i], rdiags[i] = true, true
	}
	    
	for row := 0; row<n; {
		solved = false
		for col := rows[row]; col<n; col++ {
			if cols[col] && diags[row+col] && rdiags[n-1+row-col] {
				cols[col], diags[row+col], rdiags[n-1+row-col] = false, false, false
				rows[row] = col
				if row==n-1 {
					solutions = addSolution(registerSolution(n, rows/*, cols, diags, rdiags*/), solutions)
					cols[col], diags[row+col], rdiags[n-1+row-col] = true, true, true
					rows[row] = 0
				} else {
					solved = true
					break
				}
			}
		}
		if solved {
			row++
		}else{
			for ; ; {
				row--
				if row < 0 {
					return solutions
				}
				col := rows[row]
				cols[col], diags[row+col], rdiags[n-1+row-col] = true, true, true
				for i := row+1; i < n; i++ {
					rows[i] = 0
				}
				if( rows[row] < n-1 ){
					rows[row]++; break
				}
			}
		}
	}
	return solutions
}

func uniq(sent []Solution, sol Solution) bool {
	
	rows, rowsb :=  make([] int, sol.N, sol.N), make([] int, sol.N, sol.N)
	//for y:=0; y<sol.N-1; y++ {
		//for x:=0; x<sol.N-1; x++ {
			copy(rowsb, sol.Rows)
			//rowsb[y], rowsb[x] = rowsb[x], rowsb[y] 
			for b := 0; b < sol.N; b++ {
				for a := 0; a < sol.N; a++ {
					for r, c := range rowsb { rows[(r+a)%sol.N] = c }
					for _,s :=range sent {
						if equals(rows, s.Rows) { return false }
					}
				}
				rowsb[b], rowsb[(b+1)%sol.N] = rowsb[(b+1)%sol.N], rowsb[b] 
			}
		//}
	//}
	return true
}

// Глубоко параллельный алгоритм Тимофея Большакова.
func Bolshakov(n int) []Solution{
	if n >= 4 {
		solutionsPrev := Wirth(4)
		if n == 4 {
			return solutionsPrev
		}
		for i := 5; i<=n; i++ {
			//var wg sync.WaitGroup
			start := time.Now()
			solutions := make([]Solution, 0, 4*len(solutionsPrev)*(2*i+1))
			sent := make([]Solution, 0, len(solutionsPrev))
			sChan := make(chan Solution, 4*len(solutionsPrev)*(2*i+1))
			cnt := 0
			for _, sol := range solutionsPrev {
				if uniq(sent, sol) {
					go func(s Solution){ 
						promote(s, sChan)
					}(clone(sol))
					sent = addSolution(sol, sent)
					cnt++
				}
			}
			fmt.Printf("cnt:%d\n", cnt)
			for loop := true; loop; {
				s, ok := <- sChan
				if(!ok) { break }
				if s.N > 0 { 
					solutions = addAllSolutions(s, solutions)
				} else {
					cnt--
					if cnt==0 {	close(sChan) } 
				}
			}
			dur := time.Since(start)
			sort.Sort(&solutionsSorter{sols:solutions, comp:compare})
			//printSolutions(os.Stdout, solutions)
			fmt.Printf("All (%d) solutions for %d took %s\n", len(solutions), i, dur.String())
			solutionsPrev = solutions
		} 
		return solutionsPrev
	}	
	return make([]Solution, 0, 0)
}

func printBoolArr(arr []bool){
	for i,r := range arr {
		s := 0
		if r {
			s = 1
		}
		if i == 0  { 
			fmt.Printf("%d, ", s)
		} else if i == len(arr)-1 {
			fmt.Printf("%d\n", s)
		} else {
			fmt.Printf("%d, ", s)
		}
	}
}

func set(r, c int, ns *SolutionW) bool {
	if ns.Cols[c] && ns.Diags[r+c] && ns.RDiags[ns.N-1+r-c] { 
		ns.Rows[r] = c;
		ns.Cols[c], ns.Diags[r+c], ns.RDiags[ns.N-1+r-c] = false, false, false
		return true
	}
	return false
}

func clone(s Solution) Solution{
	var s1 Solution
	s1.N = s.N
	s1.Rows = make([] int, s.N, s.N); copy(s1.Rows, s.Rows)
	return s1
}

func promote(s Solution, sChan chan Solution) {
	//fmt.Println()
	//printSolution( os.Stdout, s )
	//sols := make([]Solution, 0, s.N*(2*s.N+1))
	need, ns := true, newSolution(s.N+1)
	rows, rowsb :=  make([] int, s.N, s.N), make([] int, s.N, s.N)
	for y:=0; y<s.N-1; y++ {
		for x:=0; x<s.N-1; x++ {
			copy(rowsb, s.Rows)
			rowsb[y], rowsb[x] = rowsb[x], rowsb[y] 
			for b := 0; b < s.N; b++ {
				for a := 0; a < s.N; a++ {
					for r, c := range rowsb { rows[(r+a)%s.N] = c }
					for ir := 0; ir<=s.N; ir++ {
						for ic :=0; ic<=s.N; ic++{
							for r,c := range rows {
								nr, nc := r, c
								if nr>=ir {nr++}
								if nc>=ic {nc++}
								if(!set(nr, nc, &ns)){ need=false; break }
							}
							if(need) {
								if set(ir, ic, &ns) { 
									//fmt.Print("V")
									//sols = addSolution(ns.Solution, sols)
									sol := clone(ns.Solution)
									sChan <- sol
								}
								//printSolution( os.Stdout, ns.Solution )
							}
							reset(&ns)
							need = true
						}
					}
				}
				rowsb[b], rowsb[(b+1)%s.N] = rowsb[(b+1)%s.N], rowsb[b] 
			}
		}
	}
	/*
	for _, sol := range sols {
		sChan <- sol
	}
	*/
	var end Solution
	sChan <- end
}