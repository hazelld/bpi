package bpi

import (
//    "fmt"
//    "reflect"
)

// The purpose of this file is too allow for easy parsing of the JSON files that
// we pull from the NBA data portal. Generally, the information we want is buried
// behind a few layers of "fluff". This makes traversing down the tree very verbose
// as you have to unwrap each layer and cast into the proper type (either
// map[string]interface{} or []interface{}). Once we have the proper part of the JSON
// tree, we can easily map this too a structure using mapstructure.
// 
// For example, if you had the following json tree (taken from schedule.json)
// {
//     "league": {
//         "standard": [{
//             "watch": {
//                 "broadcast": {
//                     "video": {
// 
//                     }
//                 }
//             },
//         }, { ... }],
//         "africa" : {},
//         ...
//     }
// }
//
// And wanted to get each individual schedule, by unpeeling everything, it would
// look like this:
//
// nba_games := unwrapPath(json, []string {"league", "standard"})
// schedules := unwrapArray(nba_games, func(i int, d interface{}) interface{} {
//      sched := Schedule{}
//      mapstructure.WeakDecode(d, &sched)
//      
//      // Unwrap the video portion
//      video := unwapPath(d, []string {"watch", "broadcast", "video"})
//      mapstructure.WeakDecode(video, &sched.Video)
// })
// 

func unwrapPath (root interface{}, fields []string) interface{} {
    current, ok := root.(map[string]interface{})
    if !ok {
        return root
    }

    var last_valid_path interface{}
    for _, field := range fields {
        last_valid_path = current[field]
        current, ok = last_valid_path.(map[string]interface{})
        if !ok {
            break
        }
        last_valid_path = current
    }
    return last_valid_path
}

func unwrapArray (root interface{}, apply func(int, interface {}) interface{}) []interface{} {
    array_data, ok := root.([]interface{})
    if !ok {
        return nil
    }

    // Apply the function to each element of array, collecting the result
    var collection []interface{}
    for i, data := range array_data {
        result := apply(i, data)
        collection = append(collection, result)
    }
    return collection
}

// This function will apply a function to a map[string]interface{}, for each of
// the keys. This allows for using these keys to build a sub-structure, for 
// an example, see the Audio stream portions of the Schedule type.
func unwrapMap (root interface{}, apply func(string, interface{}) interface{}) []interface{} {
    map_data, ok := root.(map[string]interface{})
    if !ok {
        return nil
    }

    var collection []interface{}
    for key, data := range map_data {
        result := apply(key, data)
        collection = append(collection, result)
    }
    return collection
}
