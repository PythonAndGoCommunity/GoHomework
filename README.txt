INSTALLATION AND RUN INSTRUCTIONS
===============================================================================

  Linux Installation Instructions
  ------------------------------------------------------------------------------
  Open a console and type: 
    "sudo make"

  Linux Run Instructions
  ------------------------------------------------------------------------------
  Server run:
    Open a console and type: 
      "sudo docker run gohomework"
    or
    Open a console and type: 
      "sudo make run"

  Client run:
    Open a console and type:
        docker exec -i ID /usr/local/bin/client
      where ID is the ID of the running container gohomework
    or
    Open a console and type:
        docker exec -i $(docker ps -qf ancestor=gohomework --last=1) /usr/local/bin/client
      for the latest running container gohomewrok

