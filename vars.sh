#!/bin/bash

# setup for overrides in container environment prefixed by "lg_"

export USER=${lg_USER:-maxscale1}
export PASSWORD=${lg_PASSWORD:-maxscaledemo}
export HOST=${lg_HOST:-maxscale.maxscale.svc.cluster.local}
export PORT=${lg_PORT:-3306}
export DB=${lg_DB:-employees}

export RPS=${lg_RPS:-2}
export TTR=${lg_TTR:-90}
export QUERY=${lg_QUERY:-"select emp.emp_no, emp.first_name, emp.last_name, sal.salary, sal.from_date from employees emp inner join (select emp_no, MAX(salary) as salary, from_date from salaries group by emp_no) sal on (emp.emp_no = sal.emp_no) LIMIT 50;"}


