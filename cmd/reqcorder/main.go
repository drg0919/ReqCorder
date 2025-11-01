package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reqcorder/internal/diff"
	"reqcorder/internal/history"
	"reqcorder/internal/initiator"
	"reqcorder/internal/record"
	"reqcorder/internal/request"
	"reqcorder/pkg/render"
	"reqcorder/pkg/utils"
	"strconv"
)

var VERSION = "rc"

func printErrorAndExit(errStream io.Writer, err error) {
	utils.PrintError(errStream, formatErrorMessage(err))
	for key, code := range errorCodes {
		if errors.Is(err, key) {
			os.Exit(code)
		}
	}
	os.Exit(125)
}

func formatErrorMessage(err error, args ...any) error {
	for key, message := range errorMessages {
		if errors.Is(err, key) {
			if len(args) > 0 {
				return fmt.Errorf(message, args...)
			} else {
				return errors.New(message)
			}
		}
	}
	return err
}

func runDiff(outStream io.Writer, errStream io.Writer, args []string, recordStorePath string) {
	slog.Debug("Running diff command", "args", args, "recordStorePath", recordStorePath)
	var source, target string
	var inline bool
	const (
		requestType  = "requests"
		templateType = "templates"
		responseType = "responses"
	)

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			diffCommand := flag.NewFlagSet("diff", flag.ExitOnError)
			diffCommand.StringVar(&source, "source", "", "Source")
			diffCommand.StringVar(&source, "s", "", "Source (shorthand)")
			diffCommand.StringVar(&target, "target", "", "Target")
			diffCommand.StringVar(&target, "t", "", "Target (shorthand)")
			diffCommand.BoolVar(&inline, "inline", false, "Inline diff")
			diffCommand.BoolVar(&inline, "i", false, "Inline diff (shorthand)")
			diffCommand.Usage = func() {
				utils.Fprintln(errStream, "Usage of diff:\nreqcorder diff (templates|requests|responses) -s <source_identifier> -t <target_identifier> [-i|-inline] [--verbose|-v]")
				diffCommand.PrintDefaults()
			}
			diffCommand.Usage()
			return
		}
	}

	if len(args) < 1 {
		slog.Error("No diff type provided for diff command")
		printErrorAndExit(errStream, ErrorInvalidUsage)
	}
	diffType := args[0]
	validDiffTypes := map[string]bool{
		templateType: true,
		requestType:  true,
		responseType: true,
	}
	if !validDiffTypes[diffType] {
		slog.Error("Invalid diff type provided", "diffType", diffType)
		printErrorAndExit(errStream, diff.ErrorInvalidDiffType)
	}
	diffCommand := flag.NewFlagSet("diff", flag.ExitOnError)
	diffCommand.StringVar(&source, "source", "", "Source")
	diffCommand.StringVar(&source, "s", "", "Source (shorthand)")
	diffCommand.StringVar(&target, "target", "", "Target")
	diffCommand.StringVar(&target, "t", "", "Target (shorthand)")
	diffCommand.BoolVar(&inline, "inline", false, "Inline diff")
	diffCommand.BoolVar(&inline, "i", false, "Inline diff (shorthand)")
	diffCommand.Usage = func() {
		utils.Fprintln(errStream, "Usage of diff:\nreqcorder diff (templates|requests|responses) -s <source_identifier> -t <target_identifier> [-i|-inline] [--verbose|-v]")
		diffCommand.PrintDefaults()
	}
	diffCommand.Parse(args[1:])

	slog.Debug("Processing diff command", "diffType", diffType, "source", source, "target", target)
	if source == "" || target == "" {
		slog.Error("Both source and target must be provided")
		printErrorAndExit(errStream, ErrorInvalidUsage)
	}
	switch diffType {
	case templateType:
		if inline {
			slog.Debug("Running inline diff for template", "source", source, "target", target)
			err := diff.InlineDiff(outStream, recordStorePath, source, target, "template")
			if err != nil {
				slog.Error("Failed to run inline diff for template", "error", err)
				printErrorAndExit(errStream, err)
			}
		} else {
			slog.Debug("Running default diff for template", "source", source, "target", target)
			err := diff.DefaultDiff(outStream, recordStorePath, source, target, "template")
			if err != nil {
				slog.Error("Failed to run default diff for template", "error", err)
				printErrorAndExit(errStream, err)
			}
		}
	case requestType:
		if inline {
			slog.Debug("Running inline diff for request", "source", source, "target", target)
			err := diff.InlineDiff(outStream, recordStorePath, source, target, "request")
			if err != nil {
				slog.Error("Failed to run inline diff for request", "error", err)
				printErrorAndExit(errStream, err)
			}
		} else {
			slog.Debug("Running default diff for request", "source", source, "target", target)
			err := diff.DefaultDiff(outStream, recordStorePath, source, target, "request")
			if err != nil {
				slog.Error("Failed to run default diff for request", "error", err)
				printErrorAndExit(errStream, err)
			}
		}
	case responseType:
		if inline {
			slog.Debug("Running inline diff for response", "source", source, "target", target)
			err := diff.InlineDiff(outStream, recordStorePath, source, target, "response")
			if err != nil {
				slog.Error("Failed to run inline diff for response", "error", err)
				printErrorAndExit(errStream, err)
			}
		} else {
			slog.Debug("Running default diff for response", "source", source, "target", target)
			err := diff.DefaultDiff(outStream, recordStorePath, source, target, "response")
			if err != nil {
				slog.Error("Failed to run default diff for response", "error", err)
				printErrorAndExit(errStream, err)
			}
		}
	default:
		slog.Error("Invalid diff type provided", "diffType", diffType)
		printErrorAndExit(errStream, diff.ErrorInvalidDiffType)
	}
	slog.Debug("Diff command completed successfully")
}

func runShow(outStream io.Writer, errStream io.Writer, args []string, recordStorePath string) {
	slog.Debug("Running show command", "args", args, "recordStorePath", recordStorePath)
	var request, template, response string
	showCommand := flag.NewFlagSet("show", flag.ExitOnError)
	showCommand.StringVar(&request, "request", "", "Request hash")
	showCommand.StringVar(&request, "rq", "", "Request hash (shorthand)")
	showCommand.StringVar(&template, "template", "", "Template hash")
	showCommand.StringVar(&template, "tp", "", "Template hash (shorthand)")
	showCommand.StringVar(&response, "response", "", "Response ID")
	showCommand.StringVar(&response, "re", "", "Response ID (shorthand)")
	showCommand.Usage = func() {
		utils.Fprintln(errStream, "Usage of show:\nreqcorder show (-template|-tp|-request|-rq|-response|-re) <value> [--verbose|-v]")
		showCommand.PrintDefaults()
	}
	showCommand.Parse(args)
	historyStore := history.HistoryStore{
		RecordStorePath: recordStorePath,
	}
	if request != "" {
		slog.Debug("Showing request by hash", "requestHash", request)
		content, err := historyStore.GetRequestByHash(request)
		if err != nil {
			slog.Error("Failed to get request by hash", "error", err)
			printErrorAndExit(errStream, err)
		}
		utils.Fprint(outStream, content)
	} else if template != "" {
		slog.Debug("Showing template by hash", "templateHash", template)
		content, err := historyStore.GetTemplateByHash(template)
		if err != nil {
			slog.Error("Failed to get template by hash", "error", err)
			printErrorAndExit(errStream, err)
		}
		utils.Fprint(outStream, content)
	} else if response != "" {
		slog.Debug("Showing response by ID", "responseID", response)
		content, err := historyStore.GetResponseByID(response)
		if err != nil {
			slog.Error("Failed to get response by ID", "error", err)
			printErrorAndExit(errStream, err)
		}
		utils.Fprint(outStream, content)
	} else {
		slog.Error("Invalid show type - no valid parameter provided")
		printErrorAndExit(errStream, ErrorInvalidShowType)
	}
	slog.Debug("Show command completed successfully")
}

func runExec(outStream io.Writer, errStream io.Writer, args []string, recordStorePath string) {
	slog.Debug("Running exec command", "args", args, "recordStorePath", recordStorePath)
	var minimal, quiet bool
	execCommand := flag.NewFlagSet("exec", flag.ExitOnError)
	execCommand.BoolVar(&minimal, "min", false, "Only show response body and recording info on stdout")
	execCommand.BoolVar(&quiet, "quiet", false, "No output on stdout")
	execCommand.BoolVar(&minimal, "m", false, "Only show response body and recording info on stdout (shorthand)")
	execCommand.BoolVar(&quiet, "q", false, "No output on stdout (shorthand)")
	execCommand.Usage = func() {
		utils.Fprintln(errStream, "Usage of exec:\nreqcorder exec [--min|-m|--quiet|-q] <template_path> [--verbose|-v]")
		execCommand.PrintDefaults()
	}
	execCommand.Parse(args)
	if execCommand.NArg() < 1 {
		slog.Error("No template path provided for exec command")
		printErrorAndExit(errStream, ErrorInvalidUsage)
	}
	templatePath := execCommand.Args()[0]
	slog.Debug("Reading template file", "templatePath", templatePath)
	var req request.RequestObject
	err := utils.ReadYAMLFile(templatePath, &req)
	if err != nil {
		slog.Error("Failed to read template file", "templatePath", templatePath, "error", err)
		printErrorAndExit(errStream, err)
	}
	templateYaml, err := utils.ReadFile(templatePath)
	if err != nil {
		slog.Error("Failed to read template YAML", "templatePath", templatePath, "error", err)
		printErrorAndExit(errStream, err)
	}
	slog.Debug("Validating request object")
	err = req.Validate()
	if err != nil {
		slog.Error("Failed to validate request object", "error", err)
		printErrorAndExit(errStream, err)
	}
	recordStore := record.RecordStore{
		RecordStorePath: recordStorePath,
		TemplateYaml:    templateYaml,
		Request:         &req,
	}
	if !quiet && !minimal {
		utils.Fprintln(outStream, "Request Table:")
		var reqData [][]string
		reqData = append(reqData, []string{"URL", req.URL})
		reqData = append(reqData, []string{"Method", req.Method})
		reqData = append(reqData, []string{"Body preview", utils.CreatePreview(req.Body)})
		reqData = append(reqData, []string{"Authorization header", req.Auth})
		reqData = append(reqData, []string{"Timeout", req.Timeout.String()})
		render.RenderTable(outStream, []string{"Property", "Value"}, reqData...)
		utils.Fprint(outStream, "Performing request... ")
	}
	slog.Debug("Initiating HTTP request", "url", req.URL, "method", req.Method)
	res, err := initiator.InitiateRequest(&req)
	if err != nil {
		slog.Error("Failed to initiate HTTP request", "error", err)
		recordStore.Response = res
		record_err := recordStore.Record()
		if record_err != nil {
			slog.Error("Failed to record request-response cycle", "error", record_err)
		}
		if !quiet {
			utils.Fprintf(outStream, "\nResponse ID - %s\nRequest hash - %s\nTemplate hash - %s\n\n", recordStore.ResponseID, recordStore.RequestHash, recordStore.TemplateHash)
		}
		printErrorAndExit(errStream, err)
	}
	if !quiet && !minimal {
		utils.Fprintf(outStream, "Request complete (Time taken %s)\n\n", res.Timing.Total.String())
		statusStr := strconv.Itoa(int(res.StatusCode))
		if res.StatusCode < 400 {
			statusStr += " ✅"
		} else {
			statusStr += " ❌"
		}
		var resData [][]string
		utils.Fprintln(outStream, "Response Table:")
		resData = append(resData, []string{"HTTP Status Code", statusStr})
		resData = append(resData, []string{"Body preview", utils.CreatePreview(res.Body)})
		resData = append(resData, []string{"Size (Bytes)", strconv.Itoa(int(res.Size))})
		resData = append(resData, []string{"DNS lookup time", res.Timing.DNSLookup.String()})
		resData = append(resData, []string{"TCP connection time", res.Timing.TCPConnect.String()})
		resData = append(resData, []string{"TLS handshake time", res.Timing.TLSHandshake.String()})
		resData = append(resData, []string{"Time to first byte", res.Timing.FirstByte.String()})
		resData = append(resData, []string{"Total time taken ⏳", res.Timing.Total.String()})
		importantHeaders := []string{
			"Content-Type", "Location", "Cache-Control",
			"X-Ratelimit-Remaining", "X-Ratelimit-Limit",
			"X-Ratelimit-Reset",
		}

		for _, headerName := range importantHeaders {
			if value := res.Headers[headerName]; value != "" {
				resData = append(resData, []string{headerName + " header", value})
			}
		}
		render.RenderTable(outStream, []string{"Property", "Value"}, resData...)
		utils.Fprintln(outStream, "Response Body:")
	}
	if !quiet {
		body, err := utils.Prettify(res.Body)
		if err != nil {
			slog.Error("Failed to prettify response body", "error", err)
			printErrorAndExit(errStream, err)
		}
		utils.Fprintln(outStream, body)
		utils.Fprint(outStream, "\n")
	}
	if !quiet && !minimal {
		utils.Fprint(outStream, "Recording to store... ")
	}
	slog.Debug("Recording request-response cycle")
	recordStore.Response = res
	err = recordStore.Record()
	if err != nil {
		slog.Error("Failed to record request-response cycle", "error", err)
		printErrorAndExit(errStream, err)
	}
	if !quiet {
		utils.Fprintf(outStream, "Done \nResponse ID - %s\nRequest hash - %s\nTemplate hash - %s\n\n", recordStore.ResponseID, recordStore.RequestHash, recordStore.TemplateHash)
	}
	slog.Debug("Exec command completed successfully", "responseID", recordStore.ResponseID)
}

func runList(outStream io.Writer, errStream io.Writer, args []string, recordStorePath string) {
	slog.Debug("Running list command", "args", args, "recordStorePath", recordStorePath)
	var limit uint64
	var template, request string
	const (
		requestType  = "requests"
		templateType = "templates"
		responseType = "responses"
	)

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			listCommand := flag.NewFlagSet("list", flag.ExitOnError)
			listCommand.Uint64Var(&limit, "n", 10, "Limit of records to list")
			listCommand.StringVar(&request, "request", "", "Request hash filter")
			listCommand.StringVar(&request, "rq", "", "Request hash filter (shorthand)")
			listCommand.StringVar(&template, "template", "", "Template hash filter")
			listCommand.StringVar(&template, "tp", "", "Template hash filter (shorthand)")
			listCommand.Usage = func() {
				utils.Fprintln(errStream, "Usage of list:\nreqcorder list [-n] (templates|requests|responses) [-template|-tp|-request|-rq] [--verbose|-v]")
				listCommand.PrintDefaults()
			}
			listCommand.Usage()
			return
		}
	}

	if len(args) < 1 {
		slog.Error("No list type provided for list command")
		printErrorAndExit(errStream, ErrorInvalidListType)
	}
	listType := args[0]
	validListTypes := map[string]bool{
		templateType: true,
		requestType:  true,
		responseType: true,
	}
	if !validListTypes[listType] {
		slog.Error("Invalid list type provided", "listType", listType)
		printErrorAndExit(errStream, ErrorInvalidListType)
	}

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listCommand.Uint64Var(&limit, "n", 10, "Limit of records to list")
	listCommand.StringVar(&request, "request", "", "Request hash filter")
	listCommand.StringVar(&request, "rq", "", "Request hash filter (shorthand)")
	listCommand.StringVar(&template, "template", "", "Template hash filter")
	listCommand.StringVar(&template, "tp", "", "Template hash filter (shorthand)")
	listCommand.Usage = func() {
		utils.Fprintln(errStream, "Usage of list:\nreqcorder list [-n] (templates|requests|responses) [-template|-tp|-request|-rq] [--verbose|-v]")
		listCommand.PrintDefaults()
	}
	listCommand.Parse(args[1:])
	slog.Debug("Processing list command", "listType", listType, "limit", limit, "template", template, "request", request)
	historyStore := history.HistoryStore{
		RecordStorePath: recordStorePath,
	}
	switch listType {
	case responseType:
		if request == "" && template == "" {
			slog.Debug("Listing all responses sorted by timestamp", "limit", limit)
			data, err := historyStore.GetAllResponsesSorted(limit)
			if err != nil {
				slog.Error("Failed to get all responses in sorted order", "error", err)
				printErrorAndExit(errStream, err)
			}
			utils.Fprintf(outStream, "Response History (%d responses)\n", len(data))
			render.RenderTable(outStream, []string{"Response ID", "Status Code", "Total Time", "Timestamp"}, data...)
		} else if request != "" {
			slog.Debug("Listing responses by request hash", "requestHash", request, "limit", limit)
			data, err := historyStore.GetSortedResponsesByRequestHash(request, limit)
			if err != nil {
				slog.Error("Failed to get sorted responses by request hash", "error", err)
				printErrorAndExit(errStream, err)
			}
			utils.Fprintf(outStream, "Response History (%d responses)\n", len(data))
			render.RenderTable(outStream, []string{"Response ID", "Status Code", "Total Time", "Timestamp"}, data...)
		} else {
			slog.Debug("Listing responses by template hash", "templateHash", template, "limit", limit)
			data, err := historyStore.GetSortedResponsesByTemplateHash(template, limit)
			if err != nil {
				slog.Error("Failed to get sorted responses by template hash", "error", err)
				printErrorAndExit(errStream, err)
			}
			utils.Fprintf(outStream, "Response History (%d responses)\n", len(data))
			render.RenderTable(outStream, []string{"Response ID", "Status Code", "Total Time", "Timestamp"}, data...)
		}
	case requestType:
		if template == "" {
			slog.Debug("Listing all requests sorted by modification time", "limit", limit)
			data, err := historyStore.GetAllRequestsSorted(limit)
			if err != nil {
				slog.Error("Failed to get all requests sorted", "error", err)
				printErrorAndExit(errStream, err)
			}
			utils.Fprintf(outStream, "Request History (%d requests)\n", len(data))
			render.RenderTable(outStream, []string{"Request Hash", "Template Hash", "Last Modified"}, data...)
		}
	case templateType:
		slog.Debug("Listing all templates sorted by modification time", "limit", limit)
		data, err := historyStore.GetAllTemplatesSorted(limit)
		if err != nil {
			slog.Error("Failed to get all templates sorted", "error", err)
			printErrorAndExit(errStream, err)
		}
		utils.Fprintf(outStream, "Template History (%d templates)\n", len(data))
		render.RenderTable(outStream, []string{"Template Hash", "Last Modified"}, data...)
	default:
		slog.Error("Invalid list type provided", "listType", listType)
		printErrorAndExit(errStream, ErrorInvalidListType)
	}
	slog.Debug("List command completed successfully")
}

func printRootUsage(outStream io.Writer) {
	utils.Fprintln(outStream, `Usage of ReqCorder:
reqcorder <subcommand> [flags] [args]

Subcommands:
  diff     Compare two templates, requests, or responses
  show     Display a specific template, request, or response
  exec     Execute HTTP request from a template file
  list     List templates, requests, or responses in the store

Run "reqcorder <subcommand> --help" for more details.`)
}

func printVersion(outStream io.Writer) {
	utils.Fprintf(outStream, "ReqCorder version %s\n", VERSION)
}

func getBaseDir() string {
	baseDir := os.Getenv("REQCORDER_HOME")
	if baseDir == "" {
		var err error
		baseDir, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}
	return baseDir + "/.reqcorder"
}

func initLogger(verbose bool, outStream io.Writer) *os.File {
	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	baseDir := getBaseDir()
	utils.EnsureDir(baseDir + "/logs")

	logFilePath := baseDir + "/logs/app.log"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	var w io.Writer
	if verbose {
		w = io.MultiWriter(outStream, file)
	} else {
		w = io.MultiWriter(file)
	}

	logger := slog.New(slog.NewJSONHandler(w, handlerOpts))
	slog.SetDefault(logger)
	return file
}

func main() {
	verbose := false

	var subcommandArgs []string
	for _, arg := range os.Args[2:] {
		if arg != "--verbose" && arg != "-v" {
			subcommandArgs = append(subcommandArgs, arg)
		} else if arg == "--verbose" || arg == "-v" {
			verbose = true
		}
	}

	outStream := os.Stdout
	errStream := os.Stderr
	logFile := initLogger(verbose, outStream)
	defer logFile.Close()
	slog.Debug("Starting reqcorder application", "verbose", verbose, "command", os.Args)
	if len(os.Args) < 2 {
		slog.Debug("No subcommand provided, printing usage")
		printErrorAndExit(errStream, ErrorInvalidUsage)
	}
	baseDir := getBaseDir()
	recordStorePath := baseDir + "/store"
	slog.Debug("Ensuring record store directory exists", "path", recordStorePath)
	err := utils.EnsureDir(recordStorePath)
	if err != nil {
		slog.Error("Failed to create directory", "path", recordStorePath, "error", err)
		printErrorAndExit(errStream, utils.ErrorFailedToCreateDirectory)
	}
	slog.Debug("Processing subcommand", "command", os.Args[1])
	switch os.Args[1] {
	case "diff":
		slog.Debug("Running diff command")
		runDiff(outStream, errStream, subcommandArgs, recordStorePath)
	case "show":
		slog.Debug("Running show command")
		runShow(outStream, errStream, subcommandArgs, recordStorePath)
	case "exec":
		slog.Debug("Running exec command")
		runExec(outStream, errStream, subcommandArgs, recordStorePath)
	case "list":
		slog.Debug("Running list command")
		runList(outStream, errStream, subcommandArgs, recordStorePath)
	case "help", "-h", "--help":
		slog.Debug("Printing help")
		printRootUsage(outStream)
	case "version", "-v", "--version":
		slog.Debug("Printing version")
		printVersion(outStream)
	default:
		slog.Debug("Unknown command provided", "command", os.Args[1])
		printErrorAndExit(errStream, ErrorInvalidUsage)
	}
	slog.Debug("Application finished successfully")
}
