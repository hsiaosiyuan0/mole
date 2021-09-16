package assert

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/utils"
)

func Equal(t *testing.T, except, actual interface{}, msg string) {
	if !reflect.DeepEqual(except, actual) {
		t.Fatalf("%s Except: \n%v\nActual: \n%v", msg, except, actual)
	}
}

var cyan = "\033[0;36m"
var clear = "\033[0m"

func EqualString(t *testing.T, except, actual string, msg string) {
	if !reflect.DeepEqual(except, actual) {
		fmt.Printf("%sActual:%s\n%s\n", cyan, clear, actual)
		cmd := exec.Command("bash", "-c", fmt.Sprintf(`
except=$(cat <<END
%s
END
)

actual=$(cat <<END
%s
END
)

diff --color=always -u <(echo "$except") <(echo "$actual")`, except, actual))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		t.Fatalf("Failed due to strings not equal")
	}
}

func EqualJson(t *testing.T, except, actual string) {
	exceptKV, err := utils.JsonToKeyPathAndVal(except)
	if err != nil {
		t.Fatalf("deformed except %v\n", err)
	}
	actualKV, err := utils.JsonToKeyPathAndVal(actual)
	if err != nil {
		t.Fatalf("deformed actual %v\n", err)
	}
	for k, v := range exceptKV {
		exceptV := v
		actualV := actualKV[k]
		if !reflect.DeepEqual(exceptV, actualV) {
			t.Fatalf("\nUnexpected at:\n  %s\nExcept:\n  %v\nActual:\n  %v", k, exceptV, actualV)
		}
	}
}
