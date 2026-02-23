DATA_FILE = arg.csv
PID_FILE  = pids.txt

# 全てのプロセスを1秒間隔でバックグラウンド起動
run-all:
	@if [ ! -f $(DATA_FILE) ]; then echo "Error: $(DATA_FILE) not found."; exit 1; fi
	@rm -f $(PID_FILE)
	@# コメント行と空行を除外してループ
	@grep -v '^#' $(DATA_FILE) | grep -v '^$$' | while read line; do \
		UUID=$$(echo $$line | cut -d',' -f1 | xargs); \
		VALUE=$$(echo $$line | cut -d',' -f2 | xargs); \
		echo "Starting: go run . -u $$UUID -r $$VALUE"; \
		go run . -u $$UUID -r $$VALUE & \
		echo $$! >> $(PID_FILE); \
		echo "Waiting 5 second..."; \
		sleep 5; \
	done
	@echo "All processes started. PIDs are saved in $(PID_FILE)."

# 起動した全プロセスを停止
stop-all:
	@if [ -f $(PID_FILE) ]; then \
		cat $(PID_FILE) | xargs kill; \
		rm $(PID_FILE); \
		echo "Processes stopped."; \
	else \
		echo "No PID file found."; \
	fi

.PHONY: run-all stop-all