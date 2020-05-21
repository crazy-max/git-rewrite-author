package model

import "github.com/alecthomas/kong"

// Cli holds command line args, flags and cmds
type Cli struct {
	Version   kong.VersionFlag
	Repo      string `kong:"name='repo',type:'path',default='.',help='Git repository path.'"`
	LogLevel  string `kong:"name='log-level',default='info',help='Set log level.'"`
	LogCaller bool   `kong:"name='log-caller',default='false',help='Add file:line of the caller to log output.'"`
	ConfigGet struct {
	} `kong:"cmd,name='config-get',help:'Get current user name and email from Git config.'"`
	ConfigSet struct {
		Name  string `kong:"arg,required,name='name',help='Git username.'"`
		Email string `kong:"arg,required,name='email',help='Git email.'"`
	} `kong:"cmd,name='config-set',help:'Set user name and email to Git config.'"`
	List struct {
	} `kong:"cmd,name='list',help:'Display all authors/committers.'"`
	Rewrite struct {
		Old     string `kong:"arg,required,name='old',help='Current email linked to Git author to rewrite.'"`
		Correct string `kong:"arg,required,name='correct',help='New Git name and email to set.'"`
	} `kong:"cmd,name='rewrite',help:'Rewrite an author/committer in Git history.'"`
	RewriteList struct {
		File string `kong:"arg,required,name='file',type:'path',help='Authors JSON file.'"`
	} `kong:"cmd,name='rewrite-list',help:'Rewrite a list of authors/committers in Git history.'"`
}
