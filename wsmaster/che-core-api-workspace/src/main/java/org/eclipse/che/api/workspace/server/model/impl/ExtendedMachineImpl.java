/*******************************************************************************
 * Copyright (c) 2012-2016 Codenvy, S.A.
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *   Codenvy, S.A. - initial API and implementation
 *******************************************************************************/
package org.eclipse.che.api.workspace.server.model.impl;

import org.eclipse.che.api.core.model.workspace.ExtendedMachine;
import org.eclipse.che.api.core.model.workspace.ServerConf2;

import javax.persistence.CascadeType;
import javax.persistence.ElementCollection;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.JoinColumn;
import javax.persistence.OneToMany;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;

/**
 * @author Alexander Garagatyi
 */
@Entity(name = "ExternalMachine")
public class ExtendedMachineImpl implements ExtendedMachine {

    @Id
    @GeneratedValue
    private Long id;

    @ElementCollection
    private List<String> agents;

    @ElementCollection
    private Map<String, String> attributes;

    @OneToMany(cascade = CascadeType.ALL, orphanRemoval = true)
    @JoinColumn
    private Map<String, ServerConf2Impl> servers;

    public ExtendedMachineImpl() {}

    public ExtendedMachineImpl(List<String> agents,
                               Map<String, ? extends ServerConf2> servers,
                               Map<String, String> attributes) {
        if (agents != null) {
            this.agents = new ArrayList<>(agents);
        }
        if (servers != null) {
            this.servers = servers.entrySet()
                                  .stream()
                                  .collect(Collectors.toMap(Map.Entry::getKey,
                                                            entry -> new ServerConf2Impl(entry.getValue())));
        }
        if (attributes != null) {
            this.attributes = new HashMap<>(attributes);
        }
    }

    public ExtendedMachineImpl(ExtendedMachine machine) {
        this(machine.getAgents(), machine.getServers(), machine.getAttributes());
    }

    @Override
    public List<String> getAgents() {
        return agents;
    }

    public void setAgents(List<String> agents) {
        this.agents = agents;
    }

    @Override
    public Map<String, ServerConf2Impl> getServers() {
        return servers;
    }

    public void setServers(Map<String, ServerConf2Impl> servers) {
        this.servers = servers;
    }

    @Override
    public Map<String, String> getAttributes() {
        return attributes;
    }

    public void setAttributes(Map<String, String> attributes) {
        this.attributes = attributes;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof ExtendedMachineImpl)) return false;
        ExtendedMachineImpl that = (ExtendedMachineImpl)o;
        return Objects.equals(agents, that.agents) &&
               Objects.equals(servers, that.servers) &&
               Objects.equals(attributes, that.attributes);
    }

    @Override
    public int hashCode() {
        return Objects.hash(agents, servers, attributes);
    }

    @Override
    public String toString() {
        return "ExtendedMachineImpl{" +
               "agents=" + agents +
               ", servers=" + servers +
               ", attributes=" + attributes +
               '}';
    }
}
