// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// MissingAggregation is a field data based single bucket aggregation,
// that creates a bucket of all documents in the current document set context
// that are missing a field value (effectively, missing a field or having
// the configured NULL value set). This aggregator will often be used in
// conjunction with other field data bucket aggregators (such as ranges)
// to return information for all the documents that could not be placed
// in any of the other buckets due to missing field data values.
// See: http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/search-aggregations-bucket-missing-aggregation.html
type MissingAggregation struct {
    field           string
    subAggregations map[string]Aggregation
}

func NewMissingAggregation() MissingAggregation {
    a := MissingAggregation{
        subAggregations: make(map[string]Aggregation),
    }
    return a
}

func (a MissingAggregation) Field(field string) MissingAggregation {
    a.field = field
    return a
}

func (a MissingAggregation) SubAggregation(name string, subAggregation Aggregation) MissingAggregation {
    a.subAggregations[name] = subAggregation
    return a
}

func (a MissingAggregation) Source() interface{} {
    // Example:
    //	{
    //    "aggs" : {
    //      "products_without_a_price" : {
    //        "missing" : { "field" : "price" }
    //      }
    //    }
    //	}
    // This method returns only the { "missing" : { ... } } part.

    source := make(map[string]interface{})
    opts := make(map[string]interface{})
    source["missing"] = opts

    if a.field != "" {
        opts["field"] = a.field
    }

    // AggregationBuilder (SubAggregations)
    if len(a.subAggregations) > 0 {
        aggsMap := make(map[string]interface{})
        source["aggregations"] = aggsMap
        for name, aggregate := range a.subAggregations {
            aggsMap[name] = aggregate.Source()
        }
    }

    return source
}
