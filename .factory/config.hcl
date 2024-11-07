

pipeline "main" {
    filter {}
    stages = [
        { name = "build" },
        { 
            name = "test",
            depends_on = ["build"]
        }
    ]
}

stage "build" {
    run "Install golang" {
        command = ""
    }
}