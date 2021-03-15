(def cmd (nanomsg/respondent "tcp://127.0.0.1:49910"))
(def facter (script/pipeline "facter" "-j"))
(def cosmos (script/pipeline "run-cosmos" "-v"))

(def handler (fn* [msg]
     (core/case msg.cmd
                ("facts"  (script/handler "facter" "-j" msg))
                ("date"   (script/handler "date"))
                ("cosmos" (script/handler "run-cosmos" "-v" msg)))))

(cmd (tq/process handler "handler" (cmd)))
(tq/run)
