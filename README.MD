# agenc
agenc is a little tool that simplifies the use of [age](https://github.com/FiloSottile/age) for me. An encryption and decryption key (currently only SSH) can be specified in a config file. These keys are used by agenc per default that you don't have to look for the files on your terminal everytime you use encryption.
## Usage:
```commands
╭─sb@carter ~
╰─$ agenc --help
NAME:
   agenc - encrypt and decrypt using a preset key in age for convenience 
                                                                         
USAGE:                                                                   
   agenc [global options] command [command options]                      
                                                                         
COMMANDS:                                                                
   encrypt, enc, e  encrypt                                              
   decrypt, dec, d  decrypt                                              
   config, conf, c  Set config file                                      
   help, h          Shows a list of commands or help for one command     
                                                                         
GLOBAL OPTIONS:                                                          
   --help, -h  show help
```
### Disclaimer
I'm currently learning go and wanted to use this as a learning project and will try to modify it even further. It might be insecure, instable and overly complicated.