package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"text/template"
	"time"
)

type Data struct {
	Title    string
	Subject  string
	Author   string
	DueDate  string
	Tasks    []string
	FileName string
}

func escapeLaTeX(input string) string {
	forbidden := "\\&%$#_{}~^"
	for _, bad := range forbidden {
		input = strings.ReplaceAll(input, string(bad), "\\"+string(bad))
	}
	return input
}

func (d *Data) escapeLaTeX() {
	d.Title = escapeLaTeX(d.Title)
	d.Subject = escapeLaTeX(d.Subject)
	d.Author = escapeLaTeX(d.Author)
	d.DueDate = escapeLaTeX(d.DueDate)
	for i, task := range d.Tasks {
		d.Tasks[i] = escapeLaTeX(task)
	}
}

func parseArguments() (string, Data, error) {
	// Read the command line arguments
	var template string
	var data = Data{}
	var rawTasks string

	flag.StringVar(&template, "templatefile", "", "The template to be loaded")
	flag.StringVar(&data.Subject, "subject", "", "Subject of the homwork")
	flag.StringVar(&data.Title, "title", "", "Title of the homwork")
	flag.StringVar(&data.Author, "author", "", "Author of the homwork")
	flag.StringVar(&data.DueDate, "duedate", "", "Duedate of the homwork")
	flag.StringVar(&rawTasks, "tasks", "", "Coma seperated list of the tasks")
	flag.StringVar(&data.FileName, "output", "", "Output filename of the homwork")
	flag.Parse()

	if rawTasks != "" {
		data.Tasks = strings.Split(rawTasks, ",")
		for n, v := range data.Tasks {
			data.Tasks[n] = strings.TrimSpace(v)
		}
	}

	if template == "" {
		if len(flag.Args()) == 0 {
			fmt.Println("ERROR: No template specified")
			os.Exit(1)
		}
		template = flag.Args()[0]
		if template == "" {
			fmt.Println("ERROR: No template specified")
			os.Exit(1)
		}
		if path.Ext(template) == "" {
			template = template + ".tmp"
		}

		user, err := user.Current()
		if err != nil {
			return "", Data{}, err
		}
		template = path.Join(user.HomeDir, ".mkhomework", template)
	}

	// Check if the templatefile exists
	if _, err := os.Stat(template); err != nil {
		fmt.Println("ERROR: Templete does not exist")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Ask the user for the missing fields:
	fillMissingData(&data)

	return template, data, nil
}

func readField(reader *bufio.Reader, field *string, name, value string) {
	if *field != "" {
		return
	}

	fmt.Printf("%s: ", name)
	*field, _ = reader.ReadString('\n')
	*field = strings.TrimSpace(*field)

	if *field == "" && value != "" {
		*field = value
	}
}

func fillMissingData(data *Data) {
	reader := bufio.NewReader(os.Stdin)

	readField(reader, &data.Subject, "Subject", "Subject")
	readField(reader, &data.Title, "Title", "Homework")
	readField(reader, &data.Author, "Author", "Author")
	date := time.Now().Format("02.01.2006")
	readField(reader, &data.DueDate, fmt.Sprintf("Duedate [%s]", date), date)
	if len(data.Tasks) == 0 {
		var rawTasks string
		fmt.Print("Tasks: ")
		rawTasks, _ = reader.ReadString('\n')
		rawTasks = strings.TrimSpace(rawTasks)

		if rawTasks == "" {
			data.Tasks = []string{"Problem 1"}
		} else {
			data.Tasks = strings.Split(rawTasks, ",")
			for n, v := range data.Tasks {
				data.Tasks[n] = strings.TrimSpace(v)
			}
		}

	}
	readField(reader, &data.FileName, "Output File [main.tex]", "main.tex")
}

func templateName(file string) string {
	tmp := strings.Split(file, "/")
	return tmp[len(tmp)-1]
}

func main() {
	t, d, err := parseArguments()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	d.escapeLaTeX()

	tmpl, err := template.New(templateName(t)).Delims("[[", "]]").ParseFiles(t)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(d.FileName)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(file, d)
	if err != nil {
		panic(err)
	}
}
