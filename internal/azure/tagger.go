package azure

import (
	"fmt"

	"bitbucket.org/nordcloud/tagmanager/internal/azure/session"
	"bitbucket.org/nordcloud/tagmanager/internal/rules"
	tag "bitbucket.org/nordcloud/tagmanager/internal/tagger"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Found stores
type Found struct {
	Actions  []rules.ActionItem
	Resource tag.Resource
	TagRule  rules.Rule
}

type Tagger struct {
	Session   *session.AzureSession
	Found     map[string]Found
	condMap   tag.CondFuncMap
	actionMap tag.ActionFuncMap
	DryRun    bool
	Rules     rules.TagRules
}

func NewAzureTagger(ruleDef *rules.TagRules) *Tagger {
	tagger := &Tagger{}
	tagger.Rules = *ruleDef
	tagger.DryRun = *ruleDef.DryRun
	tagger.InitActionMap()
	tagger.InitCondMap()
	tagger.Found = map[string]Found{}
	return tagger
}

func (t *Tagger) InitActionMap() {
	t.actionMap = tag.ActionFuncMap{}
	t.actionMap["addTag"] = func(p map[string]string, data *tag.Resource) error {
		err := t.createOrUpdateTag(data.ID, p["tag"], p["value"])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Action addTag did not succeed for resource %s", data.ID))
		}
		return err
	}
}

func (t *Tagger) InitCondMap() {
	t.condMap = tag.CondFuncMap{}

	t.condMap["noTags"] = func(p map[string]string, data *tag.Resource) bool {
		if len(data.Tags) == 0 {
			return true
		}
		return false
	}

	t.condMap["tagEqual"] = func(p map[string]string, data *tag.Resource) bool {
		tags := data.Tags
		if len(tags) == 0 {
			return false
		}
		for k, tag := range tags {
			if p["tag"] == k && p["value"] == *tag {
				return true
			}
		}
		return false
	}

	t.condMap["tagNotEqual"] = func(p map[string]string, data *tag.Resource) bool {
		tags := data.Tags
		if len(tags) == 0 {
			return false
		}
		for k, tag := range tags {
			if p["tag"] == k && p["value"] != *tag {
				return true
			}
		}
		return false
	}

	t.condMap["tagExists"] = func(p map[string]string, data *tag.Resource) bool {
		tags := data.Tags
		if len(tags) == 0 {
			return false
		}
		if _, ok := tags[p["tag"]]; ok {
			return true
		}
		return false

	}

	t.condMap["tagNotExists"] = func(p map[string]string, data *tag.Resource) bool {
		tags := data.Tags
		if len(tags) == 0 {
			return true
		}
		if _, ok := tags[p["tag"]]; !ok {
			return true
		}
		return false
	}

	t.condMap["regionEqual"] = func(p map[string]string, data *tag.Resource) bool {
		if p["region"] == data.Region {
			return true
		}
		return false
	}

	t.condMap["regionNotEqual"] = func(p map[string]string, data *tag.Resource) bool {
		if p["region"] != data.Region {
			return true
		}
		return false
	}

	t.condMap["rgEqual"] = func(p map[string]string, data *tag.Resource) bool {
		if p["resourceGroup"] == *data.ResourceGroup {
			return true
		}
		return false
	}

	t.condMap["rgNotEqual"] = func(p map[string]string, data *tag.Resource) bool {
		if p["resourceGroup"] != *data.ResourceGroup {
			return true
		}
		return false
	}
}

func (t Tagger) ExecuteActions() error {
	for resID, found := range t.Found {
		log.Printf("🚀  Executing actions rule '%s' on %s\n", found.TagRule.Name, resID)
		for _, action := range found.Actions {
			if t.DryRun == true {
				log.Printf("  🏜 (DryRun) [%s] Action %s (%s=%s)\n", found.TagRule.Name, action.GetType(), action["tag"], action["value"])
			} else {
				log.Printf("  🚀  [%s] Action %s (%s=%s)\n", found.TagRule.Name, action.GetType(), action["tag"], action["value"])
				resource := tag.Resource{ID: resID}
				err := t.Execute(&resource, action)
				if err != nil {
					log.Errorf("Can't fire rule %s on %s\n", action.GetType(), resource.ID)
				}
			}
		}
	}
	return nil
}

//EvaluateRules iterates over all rules and resources and checks which conditions are true. Resources for which the conditions match are saved into a tagger.Found structure
func (t Tagger) EvaluateRules(resources *[]tag.Resource) error {
	var evaled bool

	for _, resource := range *resources {
		evaled = true
		log.Debugf("🔍  Checking resource: %s (%s) \n", *resource.Name, resource.ID)
		for _, y := range t.Rules.Rules {
			for _, cond := range y.Conditions {
				evaled = t.Eval(&resource, cond)
				if !evaled {
					break
				}
			}
			if evaled {
				found := Found{Actions: y.Actions, Resource: resource, TagRule: y}
				t.Found[resource.ID] = found
			}
		}
	}
	return nil
}
