#!/bin/sh

# Specail thanks for https://github.com/charmbracelet/gum

function func_exec_outcome() {
  OUTCOME_AMOUNT=$(gum input --placeholder "input outcome amount here...");
  
  if [ "$OUTCOME_AMOUNT" -gt 0 ]; then
    echo "Outcome Amount: $OUTCOME_AMOUNT";
    EXEC_COMMAND="$EXEC_COMMAND -a $OUTCOME_AMOUNT";
    
    OUTCOME_TYPE=$(gum input --placeholder "input outcome type here...");
    if [ -z "$OUTCOME_TYPE" ]; then
      OUTCOME_TYPE="unknown";
    fi
    echo "Outcome Type: $OUTCOME_TYPE";
    EXEC_COMMAND="$EXEC_COMMAND -c $OUTCOME_TYPE";
    
    OUTCOME_DATE=$(gum input --placeholder "input outcome date here...");
    if [ -z "$OUTCOME_DATE" ]; then
      OUTCOME_DATE=$(date +%Y%m%d);
    fi
    echo "Outcome Date: $OUTCOME_DATE";
    EXEC_COMMAND="$EXEC_COMMAND -b $OUTCOME_DATE";
    
    OUTCOME_REMARK=$(gum input --placeholder "input outcome remark here...");
    echo "Outcome Remark: $OUTCOME_REMARK";
    # Has some problems when there is a blankspace in remark.
    EXEC_COMMAND="$EXEC_COMMAND -d '$OUTCOME_REMARK'";
  else
    echo "Error occur...";
    exit -2;
  fi
}

function func_exec_query() {
  echo "Please choose one of the query type below:";
  QUERY_TYPE=$(gum choose "date");
  echo "Query Type: $QUERY_TYPE";
  
  if [ "$QUERY_TYPE" = "date" ]; then
    EXEC_COMMAND="$EXEC_COMMAND date";
    QUERY_DATE=$(gum input --placeholder "input query date here...");
    if [ -z "$QUERY_DATE" ]; then
      QUERY_DATE=$(date +%Y%m%d);
    fi
    
    echo "Query Date: $QUERY_DATE";
    EXEC_COMMAND="$EXEC_COMMAND -b $QUERY_DATE";
  else
    echo "Error occur...";
    exit -2;
  fi
}

function func_exec_delete() {
  echo "Please choose one of the delete type below:";
  DELETE_TYPE=$(gum choose "id");
  echo "Delete Type: $DELETE_TYPE";
  
  if [ "$DELETE_TYPE" = "id" ]; then
    EXEC_COMMAND="$EXEC_COMMAND id";
    DELETE_ID=$(gum input --placeholder "input object id here...");
    echo "Delete Id: $DELETE_ID";
    if [ -n "$DELETE_ID" ]; then
      EXEC_COMMAND="$EXEC_COMMAND -i $DELETE_ID";
    fi
  else
    echo "Error occur...";
    exit -2;
  fi
}


# Main Function
EXEC_COMMAND="";

echo "";
echo "--- Welcome to EMM-MoneyBox Interactive Mode ---";
echo "";
echo "Please choose one of the operation below:";
EXEC_TYPE=$(gum choose "outcome" "query" "delete");
echo "Operation Type: $EXEC_TYPE";

if [ "$EXEC_TYPE" = "outcome" ]; then
  EXEC_COMMAND="go run main.go cash outcome";
  func_exec_outcome;
elif [ "$EXEC_TYPE" = "query" ]; then
  EXEC_COMMAND="go run main.go cash query";
  func_exec_query;
elif [ "$EXEC_TYPE" = "delete" ]; then
  EXEC_COMMAND="go run main.go cash delete";
  func_exec_delete;
else
  echo "Error occur...";
  exit -1;
fi

if [ -z "$EXEC_COMMAND" ]; then
  # No need to execute command, just exit in normally.
  exit 0;
fi

gum confirm "Execute Command: $EXEC_COMMAND" && gum spin --spinner dot --title "Executing..." -- sleep 0.5 && $EXEC_COMMAND;
