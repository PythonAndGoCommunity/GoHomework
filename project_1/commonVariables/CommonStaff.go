package commonVariables

//KeyState - type for keys states
type KeyState int32
const (
	//Absent shows, that this key is absent
	Absent KeyState = 0
	//Present shows, that this key is present
	Present KeyState = 1
	//Ignored shows, that we don't care about state
	Ignored KeyState = 2
	//ERROR shows, that an error occurred
	ERROR KeyState = 3
)

//ArgumentFlag - flag for argument parsing, so we know, what argument we've already encountered
type ArgumentFlag int
const (
	//Nothing means no argument encountered yet
	Nothing ArgumentFlag = 0
	//Mode means mode argument encountered
	Mode ArgumentFlag = 1
	//Port means port argument encountered
	Port ArgumentFlag = 2
	//Host means host argument encountered
	Host ArgumentFlag = 3
)

//CommandFlag - flag just for convenience to remember what command we are going to process
type CommandFlag int
const (
	//GET - it's a GET command
	GET CommandFlag = 0
	//SET - it's a SET command
	SET CommandFlag = 1
	//DEL - it's a DEL command
	DEL CommandFlag = 2
)

//Answer - struct for holding database answer and key state
type Answer struct {
	Answer string
	State KeyState
}

//UsageServerErrorMessage - error message for passing invalid arguments for server
var UsageServerErrorMessage = "USAGE: server [-m/--mode MODE(disk)] [-p/--port PORT]."
//UsageClientErrorMessage - error message for passing invalid arguments for client
var UsageClientErrorMessage = "USAGE: client [-h/--host HOST] [-p/--port PORT]."
//UsageCommandsErrorMessage - error message for passing invalid command
var UsageCommandsErrorMessage = "Available commands:\n\tGET KEY\n\tSET KEY VALUE\n\tDEL KEY [KEY[...]]"
//GetErrorMessage - error message for GET
var GetErrorMessage = "USAGE: GET KEY"
//SetErrorMessage - error message for SET
var SetErrorMessage = "USAGE: SET KEY VALUE"
//ArgumentsErrorMessage - error message for passing invalid arguments as keys or values
var ArgumentsErrorMessage = "There is something wrong with arguments."